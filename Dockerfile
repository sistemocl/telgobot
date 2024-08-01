# Utiliza la imagen oficial de Go como base
FROM golang:1.20

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum y descarga las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código de la aplicación
COPY . .

# Instala chromedp y otros paquetes necesarios
RUN go get -u github.com/chromedp/chromedp
RUN go get -u gopkg.in/tucnak/telebot.v2
RUN go get -u github.com/joho/godotenv

# Compila el programa
RUN go build -o /telegram_bot

# Instala las dependencias necesarias para chromedp en modo headless
RUN apt-get update && apt-get install -y \
    libx11-xcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxi6 \
    libxtst6 \
    libnss3 \
    libglib2.0-0 \
    libnspr4 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libxrandr2 \
    libgbm1 \
    libasound2 \
    libpangocairo-1.0-0 \
    libxshmfence1 \
    libglu1-mesa \
    wget \
    unzip

# Descarga y configura Google Chrome
RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN dpkg -i google-chrome-stable_current_amd64.deb || apt-get -f install -y

# Define las variables de entorno
ENV BOT_TOKEN=""
ENV USER=""
ENV PASSWORD=""
ENV ADMIN_USER=""
ENV ADMIN_PASS=""
ENV ADMIN2=""
ENV PASS2=""

# Define el comando por defecto
CMD ["/telegram_bot"]
