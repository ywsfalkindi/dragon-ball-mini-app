# 1. نستخدم نسخة Go خفيفة للبناء
FROM golang:1.25-alpine AS builder

# نحدد مجلد العمل داخل السيرفر
WORKDIR /app

# ننسخ ملفات تعريف المكتبات أولاً (للاستفادة من الكاش)
COPY go.mod go.sum ./
RUN go mod download

# ننسخ باقي الكود
COPY . .

# نبني التطبيق ونسميه "dragon-server"
RUN go build -o dragon-server cmd/api/main.go

# 2. المرحلة النهائية: نشغل التطبيق في حاوية صغيرة جداً
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/dragon-server .
# ننسخ ملفات المايكريشن إذا كنت تحتاجها، أو نعتمد على AutoMigrate في الكود
# COPY --from=builder /app/.env .  <-- لا ننسخ ملف env، سنضعه في إعدادات الموقع

# نفتح البورت 3000
EXPOSE 3000

# أمر التشغيل
CMD ["./dragon-server"]