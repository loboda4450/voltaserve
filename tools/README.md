# Voltaserve Tools

## Getting Started

Install [Golang](https://go.dev/doc/install).

### Build and Run

Run for development:

```shell
air
```

Build binary:

```shell
go build .
```

Build Docker image:

```shell
docker build -t voltaserve/conversion .
```

## Example Requests

### Get Image Size using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "identify",
  "args": ["-format", "%w,%h", "${input}"],
  "output": true
}
```

### Convert JPEG to PNG using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "convert",
  "args": ["${input}", "${output.png}"],
  "stdout": true
}
```

### Resize an Image using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "convert",
  "args": ["-resize", "300x", "${input}", "${output.png}"],
  "stdout": true
}
```

### Generate a Thumbnail for a PDF using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `document.pdf`

`json`:

```json
{
  "bin": "convert",
  "args": ["-thumbnail", "x250", "${input}[0]", "${output.png}"],
  "stdout": true
}
```

### Convert DOCX to PDF using LibreOffice

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `document.docx`

`json`:

```json
{
  "bin": "soffice",
  "args": [
    "--headless",
    "--convert-to",
    "pdf",
    "--outdir",
    "${output.*.pdf}",
    "${input}"
  ],
  "stdout": true
}
```

### Convert PDF to Text using Poppler

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `document.pdf`

`json`:

```json
{
  "bin": "pdftotext",
  "args": ["${input}", "${output.txt}"],
  "stdout": true
}
```

### Get TSV Data From an Image Using Tesseract

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "tesseract",
  "args": ["${input}", "${output.#.tsv}", "-l", "deu", "tsv"],
  "stdout": true
}
```

### Generate PDF with OCR Text Layer From an Image Using OCRmyPDF

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "ocrmypdf",
  "args": [
    "--rotate-pages",
    "--clean",
    "--deskew",
    "--language=kor",
    "--image-dpi=300",
    "${input}",
    "${output}"
  ],
  "stdout": true
}
```

### Build Docker Images

```shell
docker build -t voltaserve/exiftool -f Dockerfile.exiftool .
```

```shell
docker build -t voltaserve/ffmpeg -f Dockerfile.ffmpeg .
```

```shell
docker build -t voltaserve/imagemagick -f Dockerfile.imagemagick .
```

```shell
docker build -t voltaserve/libreoffice -f Dockerfile.libreoffice .
```

```shell
docker build -t voltaserve/ocrmypdf -f Dockerfile.ocrmypdf .
```

```shell
docker build -t voltaserve/poppler -f Dockerfile.poppler .
```

```shell
docker build -t voltaserve/tesseract -f Dockerfile.tesseract .
```