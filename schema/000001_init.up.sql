CREATE TABLE users (
    id SERIAL PRIMARY KEY,                                  -- Идентификатор пользователя (автоматическое увеличение)
    username VARCHAR(255) UNIQUE NOT NULL,                  -- Имя пользователя (уникальное, не может быть NULL)
    email VARCHAR(255) UNIQUE NOT NULL,                     -- Email пользователя (уникальный, не может быть NULL)
    password_hash TEXT NOT NULL,                            -- Пароль пользователя (не может быть NULL)
    created_at TIMESTAMP DEFAULT NOW(),                     -- Дата создания аккаунта
    updated_at TIMESTAMP DEFAULT NOW()                      -- Дата обновления аккаунта
);

-- Создание триггера для автоматического обновления поля updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


-- Создание таблицы warehouse
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,                                              -- Идентификатор склада (автоматическое увеличение)
    name VARCHAR(255) NOT NULL,                                         -- Название склада (не может быть NULL)
    location VARCHAR(255) NOT NULL,                                     -- Локация склада (не может быть NULL)
    user_id INT NOT NULL,                                               -- Идентификатор пользователя (внешний ключ)
    created_at TIMESTAMP DEFAULT NOW(),                                 -- Дата создания
    updated_at TIMESTAMP DEFAULT NOW(),                 -- Дата обновления
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE        -- Связь с пользователем (удаление склада при удалении пользователя)
);

-- Создание триггера для автоматического обновления поля updated_at
CREATE OR REPLACE FUNCTION update_warehouses_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_warehouses_updated_at
    BEFORE UPDATE ON warehouses
    FOR EACH ROW
    EXECUTE FUNCTION update_warehouses_updated_at_column();


-- Создание таблицы products
CREATE TABLE products (
    id SERIAL PRIMARY KEY,                                                           -- Идентификатор продукта (автоматическое увеличение)
    name VARCHAR(255) NOT NULL,                                                      -- Название продукта (не может быть NULL)
    quantity INT NOT NULL CHECK (quantity >= 0),                                     -- Количество продукта (не может быть NULL)
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),                                -- Цена продукта (не может быть NULL)
    category VARCHAR(255) NOT NULL,                                                  -- Категория продукта (не может быть NULL)
    description TEXT,                                                                -- Описание продукта
    image TEXT,                                                                      -- Путь к изображению продукта
    user_id INT NOT NULL,                                                            -- Идентификатор пользователя (внешний ключ)
    warehouse_id INT NOT NULL,                                                       -- Идентификатор склада (внешний ключ)
    created_at TIMESTAMP DEFAULT NOW(),                                              -- Дата добавления
    updated_at TIMESTAMP DEFAULT NOW(),                              -- Дата обновления
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,                    -- Связь с пользователем (удаление продукта при удалении пользователя)
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE            -- Связь со складом (удаление продукта при удалении склада)
);

-- Создание триггера для автоматического обновления поля updated_at
CREATE OR REPLACE FUNCTION update_products_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_products_updated_at_column();