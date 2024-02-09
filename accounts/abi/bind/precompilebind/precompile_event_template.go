// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package precompilebind

const tmplSourcePrecompileEventGo = `
// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package {{.Package}}

import (
	"math/big"

	"github.com/ethereum/go-ethereum/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

// CUSTOM CODE STARTS HERE
// Reference imports to suppress errors from unused imports. This code and any unnecessary imports can be removed.
var (
	_ = big.NewInt
	_ = common.Big0
	_ = contract.LogGas
)

{{$structs := .Structs}}
{{$contract := .Contract}}
/* NOTE: Events can only be emitted in state-changing functions. So you cannot use events in read-only (view) functions.
Events are generally emitted at the end of a state-changing function with AddLog method of the StateDB. The AddLog method takes 4 arguments:
	1. Address of the contract that emitted the event.
	2. Topic hashes of the event.
	3. Encoded non-indexed data of the event.
	4. Block number at which the event was emitted.
The first argument is the address of the contract that emitted the event.
Topics can be at most 4 elements, the first topic is the hash of the event signature and the rest are the indexed event arguments. There can be at most 3 indexed arguments.
Topics cannot be fully unpacked into their original values since they're 32-bytes hashes.
The non-indexed arguments are encoded using the ABI encoding scheme. The non-indexed arguments can be unpacked into their original values.
Before packing the event, you need to calculate the gas cost of the event. The gas cost of an event is the base gas cost + the gas cost of the topics + the gas cost of the non-indexed data.
See Get{EvetName}EventGasCost functions for more details.
You can use the following code to emit an event in your state-changing precompile functions (generated packer might be different)):
topics, data, err := PackMyEvent(
	topic1,
	topic2,
	data1,
	data2,
)
if err != nil {
	return nil, remainingGas, err
}
accessibleState.GetStateDB().AddLog(
	ContractAddress,
	topics,
	data,
	accessibleState.GetBlockContext().Number().Uint64(),
)
*/
{{range .Contract.Events}}
	{{$event := .}}
	{{$createdDataStruct := false}}
	{{$topicCount := 0}}
	{{- range .Normalized.Inputs}}
		{{- if .Indexed}}
			{{$topicCount = add $topicCount 1}}
			{{ continue }}
		{{- end}}
		{{- if not $createdDataStruct}}
			{{$createdDataStruct = true}}
			// {{$contract.Type}}{{$event.Normalized.Name}} represents a {{$event.Normalized.Name}} non-indexed event data raised by the {{$contract.Type}} contract.
			type {{$event.Normalized.Name}}EventData struct {
		{{- end}}
		{{capitalise .Name}} {{bindtype .Type $structs}}
	{{- end}}
	{{- if $createdDataStruct}}
			}
	{{- end}}

	// Get{{.Normalized.Name}}EventGasCost returns the gas cost of the event.
	// The gas cost of an event is the base gas cost + the gas cost of the topics + the gas cost of the non-indexed data.
	// The base gas cost and the gas cost of per topics are fixed and can be found in the contract package.
	// The gas cost of the non-indexed data depends on the data type and the data size.
	func Get{{.Normalized.Name}}EventGasCost({{if $createdDataStruct}} data {{.Normalized.Name}}EventData{{end}}) uint64 {
		gas := contract.LogGas // base gas cost
		{{if $topicCount | lt 0}}
		// Add topics gas cost ({{$topicCount}} topics)
		gas += contract.LogTopicGas * {{$topicCount}}
		{{end}}

		{{range .Normalized.Inputs}}
			{{- if not .Indexed}}
				// CUSTOM CODE STARTS HERE
				// TODO: calculate gas cost for packing the data.{{decapitalise .Name}} according to the type.
				// Keep in mind that the data here will be encoded using the ABI encoding scheme.
				// So the computation cost might change according to the data type + data size and should be charged accordingly.
				// i.e gas += LogDataGas * uint64(len(data.{{decapitalise .Name}}))
				gas += contract.LogDataGas // * ...
				// CUSTOM CODE ENDS HERE
			{{- end}}
		{{- end}}

		// CUSTOM CODE STARTS HERE
		// TODO: do any additional gas cost calculation here (only if needed)
		return gas
	}

	// Pack{{.Normalized.Name}}Event packs the event into the appropriate arguments for {{.Original.Name}}.
	// It returns topic hashes and the encoded non-indexed data.
	func Pack{{.Normalized.Name}}Event({{range .Normalized.Inputs}} {{if .Indexed}}{{decapitalise .Name}} {{bindtype .Type $structs}},{{end}}{{end}}{{if $createdDataStruct}} data {{.Normalized.Name}}EventData{{end}}) ([]common.Hash, []byte, error) {
		return {{$contract.Type}}ABI.PackEvent("{{.Original.Name}}"{{range .Normalized.Inputs}},{{if .Indexed}}{{decapitalise .Name}}{{else}}data.{{capitalise .Name}}{{end}}{{end}})
	}
	{{ if $createdDataStruct }}
		// Unpack{{.Normalized.Name}}EventData attempts to unpack non-indexed [dataBytes].
		func Unpack{{.Normalized.Name}}EventData(dataBytes []byte) ({{.Normalized.Name}}EventData, error) {
			eventData := {{.Normalized.Name}}EventData{}
			err := {{$contract.Type}}ABI.UnpackIntoInterface(&eventData, "{{.Original.Name}}", dataBytes)
			return eventData, err
		}
	{{else}}
		// Unpack{{.Normalized.Name}}Event won't be generated because the event does not have any non-indexed data.
	{{end}}
{{end}}
`
