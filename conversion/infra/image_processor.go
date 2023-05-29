package infra

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
)

type ImageProcessor struct {
	cmd    *Command
	config config.Config
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		cmd:    NewCommand(),
		config: config.GetConfig(),
	}
}

func (p *ImageProcessor) Resize(inputPath string, width int, height int, outputPath string) error {
	var widthStr string
	if width == 0 {
		widthStr = ""
	} else {
		widthStr = strconv.FormatInt(int64(width), 10)
	}
	var heightStr string
	if height == 0 {
		heightStr = ""
	} else {
		heightStr = strconv.FormatInt(int64(height), 10)
	}
	size := widthStr + "x" + heightStr
	if err := p.cmd.Exec("gm", "convert", "-resize", size, inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ThumbnailImage(inputPath string, width int, height int, outputPath string) error {
	var widthStr string
	if width == 0 {
		widthStr = ""
	} else {
		widthStr = strconv.FormatInt(int64(width), 10)
	}
	var heightStr string
	if height == 0 {
		heightStr = ""
	} else {
		heightStr = strconv.FormatInt(int64(height), 10)
	}
	size := widthStr + "x" + heightStr
	if err := p.cmd.Exec("gm", "convert", "-thumbnail", size, inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) ThumbnailBase64(inputPath string, inputSize core.ImageProps) (core.Thumbnail, error) {
	if inputSize.Width > p.config.Limits.ImagePreviewMaxWidth || inputSize.Height > p.config.Limits.ImagePreviewMaxHeight {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(inputPath))
		if inputSize.Width > inputSize.Height {
			if err := p.Resize(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return core.Thumbnail{}, err
			}
		} else {
			if err := p.Resize(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return core.Thumbnail{}, err
			}
		}
		b64, err := ImageToBase64(outputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		size, err := p.Measure(outputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		return core.Thumbnail{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	} else {
		b64, err := ImageToBase64(inputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		size, err := p.Measure(inputPath)
		if err != nil {
			return core.Thumbnail{}, err
		}
		return core.Thumbnail{
			Base64: b64,
			Width:  size.Width,
			Height: size.Height,
		}, nil
	}
}

func (p *ImageProcessor) Convert(inputPath string, outputPath string) error {
	if err := p.cmd.Exec("gm", "convert", inputPath, outputPath); err != nil {
		return err
	}
	return nil
}

func (p *ImageProcessor) Measure(path string) (core.ImageProps, error) {
	res, err := p.cmd.ReadOutput("gm", "identify", "-format", "%w,%h", path)
	if err != nil {
		return core.ImageProps{}, err
	}
	values := strings.Split(res, ",")
	width, err := strconv.Atoi(helper.RemoveNonNumeric(values[0]))
	if err != nil {
		return core.ImageProps{}, err
	}
	height, err := strconv.Atoi(helper.RemoveNonNumeric(values[1]))
	if err != nil {
		return core.ImageProps{}, err
	}
	return core.ImageProps{Width: width, Height: height}, nil
}

func (p *ImageProcessor) ToBase64(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var mimeType string
	if filepath.Ext(path) == ".svg" {
		mimeType = "image/svg+xml"
	} else {
		mimeType = http.DetectContentType(b)
	}
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(b), nil
}

type ImageToDataResult struct {
	Data                []TesseractData
	NegativeConfCount   int64
	NegativeConfPercent float32
	PositiveConfCount   int64
	PositiveConfPercent float32
}

type TesseractData struct {
	BlockNum int64
	Conf     int64
	Height   int64
	Left     int64
	Level    int64
	LineNum  int64
	PageNum  int64
	ParNum   int64
	Text     string
	Top      int64
	Width    int64
	WordNum  int64
}

func (p *ImageProcessor) ImageToData(inputPath string) (ImageToDataResult, error) {
	outFile := helper.NewId()
	if err := p.cmd.Exec("tesseract", inputPath, outFile, "tsv"); err != nil {
		return ImageToDataResult{}, err
	}
	var res = ImageToDataResult{}
	outFile = outFile + ".tsv"
	f, err := os.Open(outFile)
	if err != nil {
		return ImageToDataResult{}, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return ImageToDataResult{}, err
	}
	text := string(b)
	lines := strings.Split(text, "\n")
	lines = lines[1 : len(lines)-2]
	for _, l := range lines {
		values := strings.Split(l, "\t")
		data := TesseractData{}
		data.Level, _ = strconv.ParseInt(values[0], 10, 64)
		data.PageNum, _ = strconv.ParseInt(values[1], 10, 64)
		data.BlockNum, _ = strconv.ParseInt(values[2], 10, 64)
		data.ParNum, _ = strconv.ParseInt(values[3], 10, 64)
		data.LineNum, _ = strconv.ParseInt(values[4], 10, 64)
		data.WordNum, _ = strconv.ParseInt(values[5], 10, 64)
		data.Left, _ = strconv.ParseInt(values[6], 10, 64)
		data.Top, _ = strconv.ParseInt(values[7], 10, 64)
		data.Width, _ = strconv.ParseInt(values[8], 10, 64)
		data.Height, _ = strconv.ParseInt(values[9], 10, 64)
		data.Conf, _ = strconv.ParseInt(values[10], 10, 64)
		data.Text = values[11]
		res.Data = append(res.Data, data)
	}
	for _, v := range res.Data {
		if v.Conf < 0 {
			res.NegativeConfCount++
		} else {
			res.PositiveConfCount++
		}
	}
	if len(res.Data) > 0 {
		res.NegativeConfPercent = float32((int(res.NegativeConfCount) * 100) / len(res.Data))
		res.PositiveConfPercent = float32((int(res.PositiveConfCount) * 100) / len(res.Data))
	}
	if err := os.Remove(outFile); err != nil {
		return ImageToDataResult{}, err
	}
	return res, nil
}