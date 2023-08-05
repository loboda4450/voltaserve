package runtime

import (
	"errors"
	"voltaserve/core"
	"voltaserve/identifier"
	"voltaserve/pipeline"
)

type Dispatcher struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfPipeline        core.Pipeline
	imagePipeline      core.Pipeline
	officePipeline     core.Pipeline
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		pdfPipeline:        pipeline.NewPDFPipeline(),
		imagePipeline:      pipeline.NewImagePipeline(),
		officePipeline:     pipeline.NewOfficePipeline(),
	}
}

func (d *Dispatcher) Dispatch(opts core.PipelineRunOptions) error {
	p := d.pipelineIdentifier.Identify(opts)
	if p == core.PipelinePDF {
		if err := d.pdfPipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if p == core.PipelineOffice {
		if err := d.officePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	} else if p == core.PipelineImage {
		if err := d.imagePipeline.Run(opts); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no matching pipeline found")
}
