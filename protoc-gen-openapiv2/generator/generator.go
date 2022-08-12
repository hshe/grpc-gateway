package generator

import (
	"github.com/go-kratos/grpc-gateway/v2/internal/codegenerator"
	"github.com/go-kratos/grpc-gateway/v2/internal/descriptor"
	"github.com/go-kratos/grpc-gateway/v2/protoc-gen-openapiv2/internal/genopenapi"
	"google.golang.org/protobuf/types/pluginpb"
)

// Generator is openapi v2 generator
type Generator struct {
	Reg *descriptor.Registry
}

// NewGenerator return Generator
func NewGenerator(options ...Option) *Generator {
	gen := &Generator{
		Reg: descriptor.NewRegistry(),
	}
	gen.Reg.SetUseJSONNamesForFields(true)
	gen.Reg.SetRecursiveDepth(1024)
	gen.Reg.SetMergeFileName("apidocs")
	gen.Reg.SetDisableDefaultErrors(true)
	for _, o := range options {
		o(gen)
	}
	return gen
}

// Gen generates openapi v2 json content
func (g *Generator) Gen(req *pluginpb.CodeGeneratorRequest, onlyRPC bool) (*pluginpb.CodeGeneratorResponse, error) {
	reg := g.Reg
	if reg == nil {
		reg = NewGenerator().Reg
	}
	//reg.SetGenerateRPCMethods(onlyRPC)
	reg.SetGenerateUnboundMethods(onlyRPC)
	if err := reg.SetRepeatedPathParamSeparator("csv"); err != nil {
		return nil, err
	}

	gen := genopenapi.New(reg)

	if err := genopenapi.AddErrorDefs(reg); err != nil {
		return nil, err
	}

	if err := reg.Load(req); err != nil {
		return nil, err
	}
	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			return nil, err
		}
		targets = append(targets, f)
	}

	out, err := gen.Generate(targets)
	if err != nil {
		return nil, err
	}
	return emitFiles(out), nil
}

func emitFiles(out []*descriptor.ResponseFile) *pluginpb.CodeGeneratorResponse {
	files := make([]*pluginpb.CodeGeneratorResponse_File, len(out))
	for idx, item := range out {
		files[idx] = item.CodeGeneratorResponse_File
	}
	resp := &pluginpb.CodeGeneratorResponse{File: files}
	codegenerator.SetSupportedFeaturesOnCodeGeneratorResponse(resp)
	return resp
}
