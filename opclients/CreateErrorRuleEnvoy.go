package opclients

import (
	operatorv1alpha1 "github.com/zerokdotai/zerok-operator/api/v1alpha1"
	"google.golang.org/protobuf/types/known/structpb"
	v1alpha3Spec "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetErrorRuleCrd(zerokopSpec operatorv1alpha1.ZerokopSpec) *v1alpha3.EnvoyFilter {
	envoyFilterCrd := &v1alpha3.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.istio.io/v1alpha3",
			Kind:       "EnvoyFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "error_requests",
			Namespace: "default",
		},
		Spec: v1alpha3Spec.EnvoyFilter{
			ConfigPatches: []*v1alpha3Spec.EnvoyFilter_EnvoyConfigObjectPatch{{
				ApplyTo: v1alpha3Spec.EnvoyFilter_HTTP_FILTER,
				Match: &v1alpha3Spec.EnvoyFilter_EnvoyConfigObjectMatch{
					Context: v1alpha3Spec.EnvoyFilter_ANY,
					ObjectTypes: &v1alpha3Spec.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
						Listener: &v1alpha3Spec.EnvoyFilter_ListenerMatch{
							FilterChain: &v1alpha3Spec.EnvoyFilter_ListenerMatch_FilterChainMatch{
								Filter: &v1alpha3Spec.EnvoyFilter_ListenerMatch_FilterMatch{
									Name: "envoy.filters.network.http_connection_manager",
								},
							},
						},
					},
				},
				Patch: &v1alpha3Spec.EnvoyFilter_Patch{
					Operation: v1alpha3Spec.EnvoyFilter_Patch_INSERT_BEFORE,
					Value:     GetFailedRequestValue(zerokopSpec),
				},
			}},
		},
	}
	return envoyFilterCrd
}

func CreateHeadersCrd(zerokspec operatorv1alpha1.ZerokopSpec) []*structpb.Value {
	headers_spec := zerokspec.Http_response_headers_match.Headers
	headers_envoy := []*structpb.Value{}
	for i := 0; i < len(headers_spec); i++ {
		header_spec := headers_spec[i]
		string_match := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"exact": {
					Kind: &structpb.Value_StringValue{
						StringValue: header_spec.String_match.Exact,
					},
				},
			},
		}
		header_envoy := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"name": {
					Kind: &structpb.Value_StringValue{
						StringValue: header_spec.Name,
					},
				},
				"string_match": {
					Kind: &structpb.Value_StructValue{
						StructValue: string_match,
					},
				},
				"invert": {
					Kind: &structpb.Value_BoolValue{
						BoolValue: header_spec.Invert_match,
					},
				},
			},
		}
		header_envoy_value := &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: header_envoy,
			},
		}
		headers_envoy = append(headers_envoy, header_envoy_value)
	}
	return headers_envoy
}

func GetFailedRequestValue(zerokopSpec operatorv1alpha1.ZerokopSpec) *structpb.Struct {
	headersMatch := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"headers": {
				Kind: &structpb.Value_ListValue{
					ListValue: &structpb.ListValue{
						Values: CreateHeadersCrd(zerokopSpec),
					},
				},
			},
		},
	}
	match_config := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"http_response_headers_match": {
				Kind: &structpb.Value_StructValue{
					StructValue: headersMatch,
				},
			},
		},
	}
	static_config := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"match_config": {
				Kind: &structpb.Value_StructValue{
					StructValue: match_config,
				},
			},
		},
	}
	common_config := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"static_config": {
				Kind: &structpb.Value_StructValue{
					StructValue: static_config,
				},
			},
		},
	}
	typedConfig := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"@type": {
				Kind: &structpb.Value_StringValue{
					StringValue: "type.googleapis.com/envoy.extensions.filters.http.tap.v3.Tap",
				},
			},
			"common_config": {
				Kind: &structpb.Value_StructValue{
					StructValue: common_config,
				},
			},
		},
	}
	valueStruct := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"name": {
				Kind: &structpb.Value_StringValue{
					StringValue: "envoy.filters.http.tap",
				},
			},
			"typed_config": {
				Kind: &structpb.Value_StructValue{
					StructValue: typedConfig,
				},
			},
		},
	}
	return valueStruct
}
