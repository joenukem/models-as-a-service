package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestMaasModelRefToModel_LLMISvcUsesModelRefName(t *testing.T) {
	tests := []struct {
		name              string
		metadataName      string
		modelRefName      string
		resolvedAlias     string
		expectedID        string
		expectedKind      string
		modelRefKind      string
	}{
		{
			name:          "LLMInferenceService with resolvedModelAlias uses spec.modelRef.name",
			metadataName:  "my-model",
			modelRefName:  "qwen3-8b-fp8-dynamic",
			resolvedAlias: "publishers/ai-eng-cracow/models/qwen3-8b-fp8-dynamic",
			expectedID:    "qwen3-8b-fp8-dynamic",
			expectedKind:  kindLLMISvc,
			modelRefKind:  "LLMInferenceService",
		},
		{
			name:          "LLMInferenceService without resolvedModelAlias uses spec.modelRef.name",
			metadataName:  "my-model",
			modelRefName:  "llama-7b",
			resolvedAlias: "",
			expectedID:    "llama-7b",
			expectedKind:  kindLLMISvc,
			modelRefKind:  "LLMInferenceService",
		},
		{
			name:          "LLMInferenceService without modelRef.name falls back to metadata.name",
			metadataName:  "my-model",
			modelRefName:  "",
			resolvedAlias: "publishers/ns/models/some-model",
			expectedID:    "my-model",
			expectedKind:  kindLLMISvc,
			modelRefKind:  "LLMInferenceService",
		},
		{
			name:          "empty kind defaults to LLMInferenceService",
			metadataName:  "my-model",
			modelRefName:  "bert-base",
			resolvedAlias: "publishers/ns/models/bert-base",
			expectedID:    "bert-base",
			expectedKind:  kindLLMISvc,
			modelRefKind:  "",
		},
		{
			name:          "ExternalModel uses spec.modelRef.name",
			metadataName:  "ext-ref",
			modelRefName:  "gpt-4o-external",
			resolvedAlias: "gpt-4o-external",
			expectedID:    "gpt-4o-external",
			expectedKind:  kindExternalModel,
			modelRefKind:  "ExternalModel",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "maas.opendatahub.io/v1alpha1",
				"kind":       "MaaSModelRef",
				"metadata": map[string]any{
					"name":              tc.metadataName,
					"namespace":         "test-ns",
					"creationTimestamp": metav1.Now().Format("2006-01-02T15:04:05Z"),
				},
				"spec": map[string]any{
					"modelRef": map[string]any{
						"kind": tc.modelRefKind,
						"name": tc.modelRefName,
					},
				},
				"status": map[string]any{
					"phase":    "Ready",
					"endpoint": "https://example.com",
				},
			}}

			if tc.resolvedAlias != "" {
				_ = unstructured.SetNestedField(u.Object, tc.resolvedAlias, "status", "resolvedModelAlias")
			}

			m := maasModelRefToModel(u)
			require.NotNil(t, m)
			assert.Equal(t, tc.expectedID, m.ID)
			assert.Equal(t, tc.expectedKind, m.Kind)
		})
	}
}
