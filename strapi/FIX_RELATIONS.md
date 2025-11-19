# Исправление ошибок 500 в Strapi

## Проблема
Ошибка `Cannot read properties of undefined (reading 'joinColumn')` при загрузке Content Types в админ-панели.

## Причина
1. **Конфликт имен**: В схеме `place` поле `type` (строка) конфликтовало со связью `type` (relation)
2. **Неправильная конфигурация many-to-many связей**: Strapi требует особой конфигурации для many-to-many связей

## Исправления

### 1. Переименование поля `type` в `type_name`
В схеме `place/schema.json`:
- Было: `"type": { "type": "string" }`
- Стало: `"type_name": { "type": "string" }`

### 2. Добавление связи `type` в схему `place`
Добавлена связь с `api::type.type`:
```json
"type": {
  "type": "relation",
  "relation": "manyToOne",
  "target": "api::type.type",
  "inversedBy": "places"
}
```

### 3. Исправление many-to-many связей
Для many-to-many связей в Strapi:
- На стороне `place` и `route`: убраны `mappedBy` и `inversedBy`
- На стороне `category` и `tag`: оставлены `mappedBy`

**Было:**
```json
"categories": {
  "type": "relation",
  "relation": "manyToMany",
  "target": "api::category.category",
  "mappedBy": "places"
}
```

**Стало:**
```json
"categories": {
  "type": "relation",
  "relation": "manyToMany",
  "target": "api::category.category"
}
```

### 4. Добавление связи `places` в схему `type`
В `type/schema.json` добавлена обратная связь:
```json
"places": {
  "type": "relation",
  "relation": "oneToMany",
  "target": "api::place.place",
  "mappedBy": "type"
}
```

## Текущий статус
- ✅ Исправлены конфликты имен
- ✅ Исправлены many-to-many связи
- ✅ Добавлены недостающие связи
- ⏳ Требуется перезапуск Strapi и проверка

## Следующие шаги
1. Проверить логи Strapi на наличие ошибок
2. Открыть админ-панель и проверить Content Manager
3. Если ошибки остались, проверить права доступа для Content Types

