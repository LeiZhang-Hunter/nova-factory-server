package deviceMonitorModel

type DeviceMonitorMetricReq struct {
	DeviceId   string `json:"device_id"`
	TemplateId string `json:"template_id"`
	DataId     string `json:"data_id"`
}

type T struct {
	Code int `json:"code"`
	Data struct {
		Avatar      interface{} `json:"avatar"`
		CreateDate  string      `json:"create_date"`
		CreateTime  int64       `json:"create_time"`
		Description interface{} `json:"description"`
		Dsl         struct {
			Answer     []interface{} `json:"answer"`
			Components struct {
				AnswerSweetKeysShare struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        struct {
							Content string `json:"content"`
						} `json:"output"`
						Params struct {
							DebugInputs              []interface{} `json:"debug_inputs"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   struct {
								Content string `json:"content"`
							} `json:"output"`
							OutputVarName string        `json:"output_var_name"`
							PostAnswers   []interface{} `json:"post_answers"`
							Query         []interface{} `json:"query"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Answer:SweetKeysShare"`
				BaiduReadyColtsSniff struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        interface{}   `json:"output"`
						Params        struct {
							DebugInputs              []interface{} `json:"debug_inputs"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   interface{}   `json:"output"`
							OutputVarName            string        `json:"output_var_name"`
							Query                    []struct {
								ComponentId string `json:"component_id"`
								Type        string `json:"type"`
							} `json:"query"`
							TopN int `json:"top_n"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Baidu:ReadyColtsSniff"`
				CategorizeColdJobsRemain struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        interface{}   `json:"output"`
						Params        struct {
							CategoryDescription struct {
								百度 struct {
									Description string `json:"description"`
									Examples    string `json:"examples"`
									Index       int    `json:"index"`
									To          string `json:"to"`
								} `json:"百度"`
								知识库 struct {
									Description string `json:"description"`
									Index       int    `json:"index"`
									To          string `json:"to"`
								} `json:"知识库"`
							} `json:"category_description"`
							Cite                     bool          `json:"cite"`
							DebugInputs              []interface{} `json:"debug_inputs"`
							FrequencyPenalty         float64       `json:"frequency_penalty"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							LlmEnabledTools          []interface{} `json:"llm_enabled_tools"`
							LlmId                    string        `json:"llm_id"`
							MaxTokens                int           `json:"max_tokens"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   interface{}   `json:"output"`
							OutputVarName            string        `json:"output_var_name"`
							Parameters               []interface{} `json:"parameters"`
							PresencePenalty          float64       `json:"presence_penalty"`
							Prompt                   string        `json:"prompt"`
							Query                    []struct {
								ComponentId string `json:"component_id"`
								Type        string `json:"type"`
							} `json:"query"`
							Temperature float64 `json:"temperature"`
							TopP        float64 `json:"top_p"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Categorize:ColdJobsRemain"`
				GenerateEightEggsClean struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        interface{}   `json:"output"`
						Params        struct {
							Cite                     bool          `json:"cite"`
							DebugInputs              []interface{} `json:"debug_inputs"`
							FrequencyPenalty         float64       `json:"frequency_penalty"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							LlmEnabledTools          []interface{} `json:"llm_enabled_tools"`
							LlmId                    string        `json:"llm_id"`
							MaxTokens                int           `json:"max_tokens"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   interface{}   `json:"output"`
							OutputVarName            string        `json:"output_var_name"`
							Parameters               []interface{} `json:"parameters"`
							PresencePenalty          float64       `json:"presence_penalty"`
							Prompt                   string        `json:"prompt"`
							Query                    []interface{} `json:"query"`
							Temperature              float64       `json:"temperature"`
							TopP                     float64       `json:"top_p"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Generate:EightEggsClean"`
				GenerateMetalKnivesDrop struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        interface{}   `json:"output"`
						Params        struct {
							Cite                     bool          `json:"cite"`
							DebugInputs              []interface{} `json:"debug_inputs"`
							FrequencyPenalty         float64       `json:"frequency_penalty"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							LlmEnabledTools          []interface{} `json:"llm_enabled_tools"`
							LlmId                    string        `json:"llm_id"`
							MaxTokens                int           `json:"max_tokens"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   interface{}   `json:"output"`
							OutputVarName            string        `json:"output_var_name"`
							Parameters               []interface{} `json:"parameters"`
							PresencePenalty          float64       `json:"presence_penalty"`
							Prompt                   string        `json:"prompt"`
							Query                    []interface{} `json:"query"`
							Temperature              float64       `json:"temperature"`
							TopP                     float64       `json:"top_p"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Generate:MetalKnivesDrop"`
				RetrievalEasyBearsBeam struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        interface{}   `json:"output"`
						Params        struct {
							DebugInputs              []interface{} `json:"debug_inputs"`
							EmptyResponse            string        `json:"empty_response"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							KbIds                    []string      `json:"kb_ids"`
							KbVars                   []interface{} `json:"kb_vars"`
							KeywordsSimilarityWeight float64       `json:"keywords_similarity_weight"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   interface{}   `json:"output"`
							OutputVarName            string        `json:"output_var_name"`
							Query                    []struct {
								ComponentId string `json:"component_id"`
								Type        string `json:"type"`
							} `json:"query"`
							RerankId            string  `json:"rerank_id"`
							SimilarityThreshold float64 `json:"similarity_threshold"`
							TavilyApiKey        string  `json:"tavily_api_key"`
							TopK                int     `json:"top_k"`
							TopN                int     `json:"top_n"`
							UseKg               bool    `json:"use_kg"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Retrieval:EasyBearsBeam"`
				TemplateCruelHatsDo struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        interface{}   `json:"output"`
						Params        struct {
							Content                  string        `json:"content"`
							DebugInputs              []interface{} `json:"debug_inputs"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   interface{}   `json:"output"`
							OutputVarName            string        `json:"output_var_name"`
							Parameters               []interface{} `json:"parameters"`
							Query                    []interface{} `json:"query"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Template:CruelHatsDo"`
				TemplateLazyCougarsHappen struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        interface{}   `json:"output"`
						Params        struct {
							Content                  string        `json:"content"`
							DebugInputs              []interface{} `json:"debug_inputs"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   interface{}   `json:"output"`
							OutputVarName            string        `json:"output_var_name"`
							Parameters               []interface{} `json:"parameters"`
							Query                    []interface{} `json:"query"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []string `json:"upstream"`
				} `json:"Template:LazyCougarsHappen"`
				Begin struct {
					Downstream []string `json:"downstream"`
					Obj        struct {
						ComponentName string        `json:"component_name"`
						Inputs        []interface{} `json:"inputs"`
						Output        struct {
							Content struct {
								Field1 struct {
									Content string `json:"content"`
								} `json:"0"`
							} `json:"content"`
						} `json:"output"`
						Params struct {
							DebugInputs              []interface{} `json:"debug_inputs"`
							InforVarName             string        `json:"infor_var_name"`
							Inputs                   []interface{} `json:"inputs"`
							MessageHistoryWindowSize int           `json:"message_history_window_size"`
							Output                   struct {
								Content struct {
									Field1 struct {
										Content string `json:"content"`
									} `json:"0"`
								} `json:"content"`
							} `json:"output"`
							OutputVarName string        `json:"output_var_name"`
							Prologue      string        `json:"prologue"`
							Query         []interface{} `json:"query"`
						} `json:"params"`
					} `json:"obj"`
					Upstream []interface{} `json:"upstream"`
				} `json:"begin"`
			} `json:"components"`
			EmbedId string `json:"embed_id"`
			Graph   struct {
				Edges []struct {
					Id        string `json:"id"`
					MarkerEnd string `json:"markerEnd"`
					Source    string `json:"source"`
					Style     struct {
						Stroke      string `json:"stroke"`
						StrokeWidth int    `json:"strokeWidth"`
					} `json:"style"`
					Target       string `json:"target"`
					TargetHandle string `json:"targetHandle"`
					Type         string `json:"type"`
					ZIndex       int    `json:"zIndex"`
					SourceHandle string `json:"sourceHandle,omitempty"`
					Selected     bool   `json:"selected,omitempty"`
				} `json:"edges"`
				Nodes []struct {
					Data struct {
						Label string `json:"label"`
						Name  string `json:"name"`
						Form  struct {
							Cite                     bool          `json:"cite,omitempty"`
							FrequencyPenaltyEnabled  bool          `json:"frequencyPenaltyEnabled,omitempty"`
							FrequencyPenalty         float64       `json:"frequency_penalty,omitempty"`
							LlmId                    string        `json:"llm_id,omitempty"`
							MaxTokensEnabled         bool          `json:"maxTokensEnabled,omitempty"`
							MaxTokens                int           `json:"max_tokens,omitempty"`
							MessageHistoryWindowSize int           `json:"message_history_window_size,omitempty"`
							Parameters               []interface{} `json:"parameters,omitempty"`
							PresencePenaltyEnabled   bool          `json:"presencePenaltyEnabled,omitempty"`
							PresencePenalty          float64       `json:"presence_penalty,omitempty"`
							Prompt                   string        `json:"prompt,omitempty"`
							Temperature              float64       `json:"temperature,omitempty"`
							TemperatureEnabled       bool          `json:"temperatureEnabled,omitempty"`
							TopPEnabled              bool          `json:"topPEnabled,omitempty"`
							TopP                     float64       `json:"top_p,omitempty"`
							Query                    []struct {
								ComponentId string `json:"component_id"`
								Type        string `json:"type"`
							} `json:"query,omitempty"`
							TopN                int    `json:"top_n,omitempty"`
							Content             string `json:"content,omitempty"`
							CategoryDescription struct {
								百度 struct {
									Description string `json:"description"`
									Examples    string `json:"examples"`
									Index       int    `json:"index"`
									To          string `json:"to"`
								} `json:"百度"`
								知识库 struct {
									Description string `json:"description"`
									Index       int    `json:"index"`
									To          string `json:"to"`
								} `json:"知识库"`
							} `json:"category_description,omitempty"`
							KbIds                    []string      `json:"kb_ids,omitempty"`
							KbVars                   []interface{} `json:"kb_vars,omitempty"`
							KeywordsSimilarityWeight float64       `json:"keywords_similarity_weight,omitempty"`
							SimilarityThreshold      float64       `json:"similarity_threshold,omitempty"`
							UseKg                    bool          `json:"use_kg,omitempty"`
						} `json:"form,omitempty"`
					} `json:"data"`
					Dragging bool   `json:"dragging"`
					Id       string `json:"id"`
					Measured struct {
						Height int `json:"height"`
						Width  int `json:"width"`
					} `json:"measured"`
					Position struct {
						X float64 `json:"x"`
						Y float64 `json:"y"`
					} `json:"position"`
					Selected       bool   `json:"selected"`
					SourcePosition string `json:"sourcePosition"`
					TargetPosition string `json:"targetPosition"`
					Type           string `json:"type"`
				} `json:"nodes"`
			} `json:"graph"`
			History  [][]string `json:"history"`
			Messages []struct {
				Content string `json:"content"`
				Id      string `json:"id"`
				Role    string `json:"role"`
			} `json:"messages"`
			Path      [][]string    `json:"path"`
			Reference []interface{} `json:"reference"`
		} `json:"dsl"`
		Id           string      `json:"id"`
		Nickname     string      `json:"nickname"`
		Permission   string      `json:"permission"`
		TenantAvatar interface{} `json:"tenant_avatar"`
		Title        string      `json:"title"`
		UpdateDate   string      `json:"update_date"`
		UpdateTime   int64       `json:"update_time"`
		UserId       string      `json:"user_id"`
	} `json:"data"`
	Message string `json:"message"`
}
