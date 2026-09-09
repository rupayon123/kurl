package graphql

import (
	"encoding/json"
	"testing"
)

func TestBuildPayloadIntrospect(t *testing.T) {
	payload, err := BuildPayload("", "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var p Payload
	if err := json.Unmarshal(payload, &p); err != nil {
		t.Fatalf("invalid json payload: %v", err)
	}

	if p.Query == "" {
		t.Errorf("expected non-empty introspection query")
	}
}

func TestBuildPayloadWithVariables(t *testing.T) {
	query := `query { user(id: $id) { name } }`
	vars := `{"id": "123"}`

	payload, err := BuildPayload(query, vars, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var p Payload
	if err := json.Unmarshal(payload, &p); err != nil {
		t.Fatalf("invalid json payload: %v", err)
	}

	if p.Query != query {
		t.Errorf("query mismatch")
	}
	if p.Variables["id"] != "123" {
		t.Errorf("variables mismatch")
	}
}

func TestGenerateQueryForType(t *testing.T) {
	query := GenerateQueryForType("User")
	expected := "query GetUser {\n  User {\n    id\n    name\n  }\n}"
	if query != expected {
		t.Errorf("expected %q, got %q", expected, query)
	}
}

func TestBuildPayloadPreservesNumericVariables(t *testing.T) {
	payload, err := BuildPayload(`query ($id: ID!) { user(id: $id) { name } }`, `{"id":9007199254740993,"nested":{"ratio":0.1234567890123456789}}`, false)
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Variables json.RawMessage `json:"variables"`
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		t.Fatal(err)
	}
	var variables map[string]json.RawMessage
	if err := json.Unmarshal(result.Variables, &variables); err != nil {
		t.Fatal(err)
	}
	if string(variables["id"]) != `9007199254740993` {
		t.Fatalf("numeric ID changed: %s", variables["id"])
	}
	if string(variables["nested"]) != `{"ratio":0.1234567890123456789}` {
		t.Fatalf("decimal changed: %s", variables["nested"])
	}
}

func TestBuildPayloadRejectsInvalidVariableObjects(t *testing.T) {
	for _, variables := range []string{`{"x":1} {}`, `[]`, `1`, `{`} {
		if _, err := BuildPayload(`{ user { name } }`, variables, false); err == nil {
			t.Errorf("accepted invalid variables %s", variables)
		}
	}
}
