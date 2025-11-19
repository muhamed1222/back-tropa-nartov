#!/bin/bash

# Создаем директорию, если её нет
mkdir -p "/Users/kelemetovmuhamed/Documents/тропа нартов /front"

# Переходим в директорию
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /front"

# Клонируем репозиторий с веткой final_back
git clone --branch final_back --single-branch https://github.com/muhamed1222/back.git

# Если клонирование прошло успешно, копируем содержимое папки Downloads/back-now3
if [ -d "back/Downloads/back-now3" ]; then
    echo "Копирую содержимое папки Downloads/back-now3..."
    cp -r back/Downloads/back-now3/* . 2>/dev/null || cp -r back/Downloads/back-now3/. . 2>/dev/null
    echo "Готово!"
else
    echo "Папка Downloads/back-now3 не найдена в репозитории"
fi

