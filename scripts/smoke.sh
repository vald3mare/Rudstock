#!/usr/bin/env bash
# Смоук-тест реализованных ручек: на каждый запрос сверяет код ответа
# и подстроку тела (обычно slug ошибки), в конце печатает итог.
#
# Нужен запущенный сервер с базой после миграций. Создаёт свою категорию
# smoke-<timestamp> и две карточки в ней. В конце удаляет их через DELETE-ручки
# (это заодно и проверка удаления), а cleanup в trap добирает остатки,
# если скрипт прервали (Ctrl+C) раньше.
#
# Запуск:
#   scripts/smoke.sh                     # сервер на localhost:8081
#   API=localhost:8082 scripts/smoke.sh  # другой адрес
#
# Код выхода: 0 если всё прошло, 1 если хоть одна проверка упала.

API=${API:-localhost:8081}
J='Content-Type: application/json'
pass=0
fail=0

# ждём сервер до ~10 секунд, чтобы скрипт можно было запускать сразу после старта
for _ in $(seq 1 30); do
  curl -s -o /dev/null "$API/health" && break
  sleep 0.3
done

# check <название> <ожидаемый код> <подстрока тела или -> <аргументы curl...>
# Тело последнего ответа остаётся в LAST_BODY, из него достаём созданные id.
check() {
  local name=$1 want=$2 sub=$3
  shift 3
  local out code body
  out=$(curl -s -w $'\n%{http_code}' "$@")
  code=${out##*$'\n'}
  body=${out%$'\n'*}
  if [[ $code == "$want" && ($sub == - || $body == *"$sub"*) ]]; then
    pass=$((pass + 1))
    printf 'OK   %-3s %s\n' "$code" "$name"
  else
    fail=$((fail + 1))
    printf 'FAIL %-3s %s (want %s %s)\n     %s\n' "$code" "$name" "$want" "$sub" "$body"
  fi
  LAST_BODY=$body
}

# уникальное имя: иначе повторный запуск упрётся в 409 на создании категории
NAME="smoke-$(date +%s%N)"

# id созданного, заполняются по ходу; пустые значит "не создали"
CAT=""
CARD=""
CARD2=""

# cleanup молча удаляет всё созданное: сначала карточки, потом категорию,
# иначе база не даст удалить категорию (409). Уже удалённое вернёт 404, это нормально.
cleanup() {
  [[ -n $CARD ]] && curl -s -o /dev/null -X DELETE "$API/admin/cards/$CARD"
  [[ -n $CARD2 ]] && curl -s -o /dev/null -X DELETE "$API/admin/cards/$CARD2"
  [[ -n $CAT ]] && curl -s -o /dev/null -X DELETE "$API/admin/categories/$CAT"
}
trap cleanup EXIT

echo "== health ($API)"
check "health" 200 '"ok"' "$API/health"

echo "== POST /admin/categories"
check "create category" 201 category_id -X POST "$API/admin/categories" -H "$J" -d "{\"name\":\"$NAME\"}"
CAT=$(sed -E 's/.*"category_id":([0-9]+).*/\1/' <<<"$LAST_BODY")
check "duplicate" 409 category-exists -X POST "$API/admin/categories" -H "$J" -d "{\"name\":\"$NAME\"}"
check "duplicate with spaces (trim)" 409 category-exists -X POST "$API/admin/categories" -H "$J" -d "{\"name\":\"  $NAME  \"}"
check "empty name" 400 invalid-category-name -X POST "$API/admin/categories" -H "$J" -d '{"name":""}'
check "spaces name" 400 invalid-category-name -X POST "$API/admin/categories" -H "$J" -d '{"name":"   "}'
check "broken json" 400 invalid-json-body -X POST "$API/admin/categories" -H "$J" -d '{"name":'
check "unknown field" 400 invalid-json-body -X POST "$API/admin/categories" -H "$J" -d '{"title":"x"}'

echo "== GET /categories"
check "list categories" 200 "$NAME" "$API/categories"

echo "== POST /admin/cards"
check "create card" 201 card_id -X POST "$API/admin/cards" -H "$J" \
  -d "{\"category_id\":$CAT,\"description\":\"Шина летняя 205/55 R16\",\"price\":14999,\"photo_url\":\"s3://photos/photo_1.jpg\"}"
CARD=$(sed -E 's/.*"card_id":"([^"]+)".*/\1/' <<<"$LAST_BODY")
check "required fields only" 201 card_id -X POST "$API/admin/cards" -H "$J" -d "{\"category_id\":$CAT,\"price\":100}"
CARD2=$(sed -E 's/.*"card_id":"([^"]+)".*/\1/' <<<"$LAST_BODY")
check "price 0" 400 invalid-price -X POST "$API/admin/cards" -H "$J" -d "{\"category_id\":$CAT,\"price\":0}"
check "no category_id" 400 invalid-category-id -X POST "$API/admin/cards" -H "$J" -d '{"price":100}'
check "unknown category" 404 category-not-found -X POST "$API/admin/cards" -H "$J" -d '{"category_id":999999,"price":100}'

echo "== GET /cards"
# в свежей категории ровно 2 карточки, созданные выше
check "list all" 200 '"total"' "$API/cards"
check "by category" 200 '"total":2' "$API/cards?category_id=$CAT&page=1&limit=20"
check "limit=1 keeps total" 200 '"total":2' "$API/cards?category_id=$CAT&limit=1"
check "page past end" 200 '"items":[],"total":2' "$API/cards?category_id=$CAT&page=1000"
check "limit 1000 clamps" 200 '"items"' "$API/cards?limit=1000"
check "unknown category" 200 '"items":[],"total":0' "$API/cards?category_id=999999"
check "limit not int" 400 invalid-query-param "$API/cards?limit=abc"
check "negative page" 400 invalid-pagination "$API/cards?page=-1"
check "negative category" 400 invalid-category-id "$API/cards?category_id=-1"

echo "== GET /cards/{id}"
check "get card" 200 "\"id\":\"$CARD\"" "$API/cards/$CARD"
check "not uuid" 400 invalid-path-param "$API/cards/abc"
check "not found" 404 card-not-found "$API/cards/00000000-0000-0000-0000-000000000000"

echo "== PATCH /admin/cards/{id}"
# после изменений перечитываем карточку через GET: проверяем, что сохранилось, а не только ответ PATCH
check "price+description" 200 '"price":12999' -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{"price":12999,"description":"Распродажа"}'
check "  description changed" 200 '"description":"Распродажа"' "$API/cards/$CARD"
check "  photo untouched" 200 '"photo_url":"s3://photos/photo_1.jpg"' "$API/cards/$CARD"
check "clear photo" 200 '"photo_url":""' -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{"photo_url":""}'
check "null keeps field" 200 '"description":"Распродажа"' -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{"description":null,"price":13999}'
check "  price changed" 200 '"price":13999' "$API/cards/$CARD"
check "empty patch" 400 empty-patch -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{}'
check "price 0" 400 invalid-price -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{"price":0}'
check "category 0" 400 invalid-category-id -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{"category_id":0}'
check "unknown category" 404 category-not-found -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{"category_id":999999}'
check "unknown field" 400 invalid-json-body -X PATCH "$API/admin/cards/$CARD" -H "$J" -d '{"title":"x"}'
check "not uuid" 400 invalid-path-param -X PATCH "$API/admin/cards/abc" -H "$J" -d '{"price":100}'
check "not found" 404 card-not-found -X PATCH "$API/admin/cards/00000000-0000-0000-0000-000000000000" -H "$J" -d '{"price":100}'

echo "== routing"
check "wrong method" 405 - -X PUT "$API/cards/$CARD"

# удаление в конце: это и проверка DELETE-ручек, и уборка за собой
echo "== DELETE /admin/categories/{id} (with cards)"
check "category has cards" 409 category-has-cards -X DELETE "$API/admin/categories/$CAT"
check "  category still exists" 200 "$NAME" "$API/categories"

echo "== DELETE /admin/cards/{id}"
check "delete card" 204 - -X DELETE "$API/admin/cards/$CARD"
check "  card gone" 404 card-not-found "$API/cards/$CARD"
check "delete again" 404 card-not-found -X DELETE "$API/admin/cards/$CARD"
check "delete second card" 204 - -X DELETE "$API/admin/cards/$CARD2"
check "not uuid" 400 invalid-path-param -X DELETE "$API/admin/cards/abc"

echo "== DELETE /admin/categories/{id}"
check "delete empty category" 204 - -X DELETE "$API/admin/categories/$CAT"
check "delete again" 404 category-not-found -X DELETE "$API/admin/categories/$CAT"
check "not int" 400 invalid-path-param -X DELETE "$API/admin/categories/abc"

echo
echo "category=$CAT card=$CARD card2=$CARD2"
echo "passed=$pass failed=$fail"

[[ $fail -eq 0 ]]