-- 1. جدول المستخدمين (The Fighters)
CREATE TABLE users (
    -- SERIAL: يعني رقم يزيد تلقائياً (1, 2, 3...)
    -- PRIMARY KEY: هذا هو الرقم المميز الذي لا يتكرر أبداً
    id SERIAL PRIMARY KEY,
    
    -- BIGINT: لأن رقم تليجرام ضخم جداً لا يكفيه الـ Int العادي
    -- UNIQUE: ممنوع تكرار نفس رقم التليجرام (المستخدم يسجل مرة واحدة)
    telegram_id BIGINT UNIQUE NOT NULL,
    
    -- VARCHAR(255): نص طوله أقصى حد 255 حرف
    username VARCHAR(255),
    
    -- DEFAULT 10: أي لاعب جديد يبدأ بـ 10 طاقة فوراً
    energy INT DEFAULT 10,
    
    -- TIMESTAMP: يسجل التاريخ والوقت
    -- DEFAULT NOW(): يضع وقت التسجيل الحالي تلقائياً
    created_at TIMESTAMP DEFAULT NOW()
);

-- 2. جدول الأسئلة (The Scrolls of Wisdom)
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    question_text TEXT NOT NULL, -- TEXT: نص طويل جداً
    option_a VARCHAR(255) NOT NULL,
    option_b VARCHAR(255) NOT NULL,
    option_c VARCHAR(255) NOT NULL,
    option_d VARCHAR(255) NOT NULL,
    correct_option CHAR(1) NOT NULL, -- حرف واحد فقط (A, B, C, D)
    difficulty INT DEFAULT 1 -- 1: Easy, 2: Medium, 3: Hard
);

-- 3. جدول النتائج (The Battle History)
CREATE TABLE scores (
    id SERIAL PRIMARY KEY,
    
    -- هذا هو "المفتاح الأجنبي" (Foreign Key)
    -- REFERENCES users(id): هذا العمود مرتبط بعمود الـ id في جدول users
    -- ON DELETE CASCADE: لو حذفنا المستخدم، تحذف كل نتائجه تلقائياً (تنظيف)
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    
    points INT NOT NULL,
    obtained_at TIMESTAMP DEFAULT NOW()
);