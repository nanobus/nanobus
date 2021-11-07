package spec_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nanobus/nanobus/spec"
)

func TestRegistry(t *testing.T) {
	r := spec.Registry{}

	loader := func(config interface{}) ([]*spec.Namespace, error) {
		return nil, nil
	}
	namedLoader := func() (string, spec.Loader) {
		return "test", loader
	}

	r.Register(namedLoader)

	assert.Equal(t, fmt.Sprintf("%v", spec.Loader(loader)), fmt.Sprintf("%p", r["test"]))
}

func TestNamespace(t *testing.T) {
	expectedBytes, err := os.ReadFile("testdata/expected.json")
	if err != nil {
		t.FailNow()
	}

	arg1 := spec.NewArgument("arg1", "val1")
	assert.Equal(t, "val1", arg1.ValueString())

	anno := spec.NewAnnotation("test").
		AddArguments(
			arg1,
			spec.NewArgument("arg2", "val2"))

	annoMap := anno.ToMap()
	assert.Equal(t, map[string]interface{}{
		"arg1": "val1",
		"arg2": "val2",
	}, annoMap)

	var argStruct struct {
		Arg1 string `mapstructure:"arg1"`
		Arg2 string `mapstructure:"arg2"`
	}
	if assert.NoError(t, anno.ToStruct(&argStruct)) {
		assert.Equal(t, "val1", argStruct.Arg1)
		assert.Equal(t, "val2", argStruct.Arg2)
	}

	arg1Argument, ok := anno.Argument("arg1")
	if assert.True(t, ok) {
		assert.Same(t, arg1, arg1Argument)
	}

	dogBreed := spec.NewEnum("DogBreed", "Enumeration for the type of animal").
		AddAnnotations(anno).
		AddAnnotation(anno). // Should not be added twice
		AddValues(
			spec.NewEnumValue("corgi", "Pembroke Welsh Corgi", "Pembroke Welsh Corgi", 1).
				AddAnnotations(anno).
				AddAnnotation(anno), // Should not be added twice
			spec.NewEnumValue("golden_retriever", "Golden Retriever", "Golden Retriever", 2).
				AddAnnotations(anno))

	name := spec.NewField("name", "The pet's name", &spec.TypeRef{
		Kind: spec.KindString,
	}, nil).
		AddAnnotations(anno).
		AddAnnotation(anno)

	dog := spec.NewType("Dog", "A dog").
		AddAnnotations(anno).
		AddAnnotation(anno). // Should not be added twice
		AddFields(
			name,
			spec.NewField("breed", "The dog's breed", &spec.TypeRef{
				Kind: spec.KindEnum,
				Enum: dogBreed,
			}, nil).
				AddAnnotations(anno),
			spec.NewField("parentDogIds", "The dog's parents", &spec.TypeRef{
				Kind: spec.KindList,
				ListType: &spec.TypeRef{
					Kind: spec.KindString,
				},
			}, nil),
			spec.NewField("birthDate", "The dog's birth date", &spec.TypeRef{
				Kind: spec.KindOptional,
				OptionalType: &spec.TypeRef{
					Kind: spec.KindDateTime,
				},
			}, nil),
			spec.NewField("diet", "The dog's food intake", &spec.TypeRef{
				Kind: spec.KindMap,
				MapKeyType: &spec.TypeRef{
					Kind: spec.KindString,
				},
				MapValueType: &spec.TypeRef{
					Kind: spec.KindString,
				},
			}, nil),
		)

	nameField, ok := dog.Field("name")
	if assert.True(t, ok) {
		assert.Same(t, name, nameField)
	}

	cat := spec.NewType("Cat", "A cat").
		AddAnnotations(anno).
		AddAnnotation(anno). // Should not be added twice
		AddFields(
			name,
			spec.NewField("liveRemaining", "How many lives are remaiing", &spec.TypeRef{
				Kind: spec.KindU8,
			}, nil),
		)

	animal := spec.NewUnion("Animal", "A union of animal types").
		AddTypes(&spec.TypeRef{
			Kind: spec.KindType,
			Type: dog,
		}).
		AddType(&spec.TypeRef{
			Kind: spec.KindType,
			Type: cat,
		}).
		AddAnnotations(anno).
		AddAnnotation(anno)

	response := spec.NewType("Response", "Greeting response").
		AddFields(
			spec.NewField("message", "The greeting message", &spec.TypeRef{
				Kind: spec.KindString,
			}, nil),
		).AddAnnotations(anno).
		AddAnnotations(anno)

	sayHello := spec.NewOperation("sayHello", "Say hello", true,
		spec.NewType("sayHelloArgs", "arguments for sayHello").
			AddFields(
				spec.NewField("name", "Name of the person to greet", &spec.TypeRef{
					Kind: spec.KindString,
				}, "World"),
			),
		&spec.TypeRef{
			Kind: spec.KindType,
			Type: response,
		}).
		AddAnnotations(anno).
		AddAnnotation(anno)

	getAnimal := spec.NewOperation("getAnimal", "Retrieve an animal", false,
		spec.NewType("getAnimalArgs", "arguments for sayHello").
			AddFields(
				spec.NewField("animalId", "ID of the animal", &spec.TypeRef{
					Kind: spec.KindString,
				}, nil),
			),
		&spec.TypeRef{
			Kind:  spec.KindUnion,
			Union: animal,
		}).
		AddAnnotations(anno).
		AddAnnotation(anno)

	service := spec.NewService("Hello", "Greetings").
		AddAnnotations(
			spec.NewAnnotation("service")).
		AddAnnotation(anno).
		AddOperations(sayHello, getAnimal)

	testAnno, ok := service.Annotation("test")
	if assert.True(t, ok) {
		assert.Same(t, anno, testAnno)
	}

	oper, ok := service.Operation("sayHello")
	if assert.True(t, ok) {
		assert.Same(t, sayHello, oper)
	}

	ns := spec.NewNamespace("greetings.v1").
		AddAnnotations(
			spec.NewAnnotation("anno").AddArguments(
				spec.NewArgument("arg1", "val1"),
				spec.NewArgument("arg2", "val2")),
		).
		AddAnnotation(anno).
		AddServices(service).
		AddTypes(response, dog, cat).
		AddUnions(animal).
		AddEnums(dogBreed)

	dogType, ok := ns.Type("Dog")
	if assert.True(t, ok) {
		assert.Same(t, dog, dogType)
	}

	animalUnion, ok := ns.Union("Animal")
	if assert.True(t, ok) {
		assert.Same(t, animal, animalUnion)
	}

	breedEnum, ok := ns.Enum("DogBreed")
	if assert.True(t, ok) {
		assert.Same(t, dogBreed, breedEnum)
	}

	namespaces := spec.Namespaces{}.AddNamespaces(ns)

	oper, ok = namespaces.Operation("greetings.v1", "Hello", "sayHello")
	if assert.True(t, ok) {
		assert.Same(t, sayHello, oper)
	}

	serv, ok := ns.Service("Hello")
	if assert.True(t, ok) {
		assert.Same(t, service, serv)
	}

	actualBytes, err := json.MarshalIndent(ns, "", "  ")
	require.NoError(t, err)
	fmt.Println(string(actualBytes))

	var expected, actual interface{}
	require.NoError(t, json.Unmarshal(expectedBytes, &expected))
	require.NoError(t, json.Unmarshal(actualBytes, &actual))

	assert.Equal(t, expected, actual)
}
