// Copyright © 2023 Horizoncd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"os"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/stretchr/testify/assert"
)

func TestTextTemplateFromTextInCode(t *testing.T) {
	textNotTemplate := `
{
    "type":"object",
    "properties":{
        "app":{
            "title":"",
            "type":"object",
            "properties":{
                "spec":{
                    "type":"object",
                    "title":"规格",
                    "properties":{
                        "replicas":{
                            "type":"integer",
                            "title":"副本数",
                            "default": 1,
                            "minimum":0,
                            "maximum":50,
                        },
                        "resource":{
                            "type":"string",
                            "title":"规格",
                            "default":"x-small",
                            "oneOf":[
                                {
                                    "enum":[
                                        "x-small"
                                    ],
                                    "title":"x-small(1C2G)"
                                },
                                {
                                    "enum":[
                                        "small"
                                    ],
                                    "title":"small(2C4G)"
                                },
                                {
                                    "enum":[
                                        "middle"
                                    ],
                                    "title":"middle(4C8G)"
                                },
                                {
                                    "enum":[
                                        "large"
                                    ],
                                    "title":"large(8C16G)"
                                },
                                {
                                    "enum":[
                                        "x-large"
                                    ],
                                    "title":"x-large(16C32G)"
                                }
                            ]
                        }
                    }
                },
                "strategy":{
                    "type":"object",
                    "title":"发布策略",
                    "properties":{
                        "stepsTotal":{
                            "type":"integer",
                            "title":"发布批次（多批次情况下，第一批默认为1个实例）",
                            "default": 1,
                            "enum":[
                                1,
                                2,
                                3,
                                4,
                                5
                            ]
                        },
                        "pauseType":{
                            "type":"string",
                            "title":"暂停策略",
                            "default": "all",
                            "oneOf":[
                                {
                                    "enum":[
                                        "first"
                                    ],
                                    "title":"第一批暂停"
                                },
                                {
                                    "enum":[
                                        "all"
                                    ],
                                    "title":"全部暂停"
                                },
                                {
                                    "enum":[
                                        "none"
                                    ],
                                    "title":"全不暂停"
                                }
                            ]
                        }
                    }
                },
                "params":{
                    "title":"参数",
                    "type":"object",
                    "properties":{
                        "mainClassName":{
                            "type":"string"
                        },
                        "xmx":{
                            "type":"string",
                            "default":"512",
                            "pattern":"^\\d*$"
                        },
                        "xms":{
                            "type":"string",
                            "default":"512",
                            "pattern":"^\\d*$"
                        },
                        "xdebugAddress":{
                            "type":"string",
                            "pattern":"^\\d*$"
                        },
                        "jvmExtra":{
                            "type":"string"
                        }
                    },
                    "required":[
                        "mainClassName"
                    ]
                },
                "health":{
                    "title":"健康检查",
                    "type":"object",
                    "properties":{
                        "port":{
                            "type":"integer",
                            "minimum":1024,
                            "maximum":65535
                        },
                        "check":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"存活状态",
                            "description":"存活状态会在应用运行期间检测应用健康情况，检测失败时会对应用进行重启。接口如: /api/test"
                        },
                        "status":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"就绪状态",
                            "description":"就绪状态会在应用运行期间检测应用是否处于上线状态，检测失败时显示下线状态。接口如: /health/status"
                        },
                        "online":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"上线",
                            "description":"上线接口会在应用启动之后进行调用，如果调用失败，则应用启动失败。接口如: /health/online"
                        },
                        "offline":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"下线",
                            "description":"下线接口会在应用停止之前进行调用，如果调用失败，则忽略。接口如: /health/offline"
                        }
                    },
                    "dependencies":{
                        "check":[
                            "port"
                        ],
                        "status":[
                            "port"
                        ],
                        "online":[
                            "port"
                        ],
                        "offline":[
                            "port"
                        ]
                    }
                }
            }
        }
    }
}
`
	template1 := template.Must(template.New("template1.text").Funcs(sprig.TxtFuncMap()).Parse(textNotTemplate))

	vals := make(map[interface{}]interface{})
	vals["TomAddInt"] = "100"
	err := template1.ExecuteTemplate(os.Stdout, "template1.text", vals)
	assert.Nil(t, err)

	vals["TomAddInt"] = "102"

	textTemplate := `
{
    "type":"object",
    "properties":{
        "app":{
            "title":"",
            "type":"object",
            "properties":{
                "spec":{
                    "type":"object",
                    "title":"规格",
                    "properties":{
                        "replicas":{
                            "type":"integer",
                            "title":"副本数",
                            "default": 1,
                            "minimum":0,
                            "maximum":{{ add .TomAddInt 2 }}
                        },
                        "resource":{
                            "type":"string",
                            "title":"规格",
                            "default":"x-small",
                            "oneOf":[
                                {
                                    "enum":[
                                        "x-small"
                                    ],
                                    "title":"x-small(1C2G)"
                                },
                                {
                                    "enum":[
                                        "small"
                                    ],
                                    "title":"small(2C4G)"
                                },
                                {
                                    "enum":[
                                        "middle"
                                    ],
                                    "title":"middle(4C8G)"
                                },
                                {
                                    "enum":[
                                        "large"
                                    ],
                                    "title":"large(8C16G)"
                                },
                                {
                                    "enum":[
                                        "x-large"
                                    ],
                                    "title":"x-large(16C32G)"
                                }
                            ]
                        }
                    }
                },
                "strategy":{
                    "type":"object",
                    "title":"发布策略",
                    "properties":{
                        "stepsTotal":{
                            "type":"integer",
                            "title":"发布批次（多批次情况下，第一批默认为1个实例）",
                            "default": 1,
                            "enum":[
                                1,
                                2,
                                3,
                                4,
                                5
                            ]
                        },
                        "pauseType":{
                            "type":"string",
                            "title":"暂停策略",
                            "default": "all",
                            "oneOf":[
                                {
                                    "enum":[
                                        "first"
                                    ],
                                    "title":"第一批暂停"
                                },
                                {
                                    "enum":[
                                        "all"
                                    ],
                                    "title":"全部暂停"
                                },
                                {
                                    "enum":[
                                        "none"
                                    ],
                                    "title":"全不暂停"
                                }
                            ]
                        }
                    }
                },
                "params":{
                    "title":"参数",
                    "type":"object",
                    "properties":{
                        "mainClassName":{
                            "type":"string"
                        },
                        "xmx":{
                            "type":"string",
                            "default":"512",
                            "pattern":"^\\d*$"
                        },
                        "xms":{
                            "type":"string",
                            "default":"512",
                            "pattern":"^\\d*$"
                        },
                        "xdebugAddress":{
                            "type":"string",
                            "pattern":"^\\d*$"
                        },
                        "jvmExtra":{
                            "type":"string"
                        }
                    },
                    "required":[
                        "mainClassName"
                    ]
                },
                "health":{
                    "title":"健康检查",
                    "type":"object",
                    "properties":{
                        "port":{
                            "type":"integer",
                            "minimum":1024,
                            "maximum":65535
                        },
                        "check":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"存活状态",
                            "description":"存活状态会在应用运行期间检测应用健康情况，检测失败时会对应用进行重启。接口如: /api/test"
                        },
                        "status":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"就绪状态",
                            "description":"就绪状态会在应用运行期间检测应用是否处于上线状态，检测失败时显示下线状态。接口如: /health/status"
                        },
                        "online":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"上线",
                            "description":"上线接口会在应用启动之后进行调用，如果调用失败，则应用启动失败。接口如: /health/online"
                        },
                        "offline":{
                            "type":"string",
                            "pattern":"^/.*$",
                            "title":"下线",
                            "description":"下线接口会在应用停止之前进行调用，如果调用失败，则忽略。接口如: /health/offline"
                        }
                    },
                    "dependencies":{
                        "check":[
                            "port"
                        ],
                        "status":[
                            "port"
                        ],
                        "online":[
                            "port"
                        ],
                        "offline":[
                            "port"
                        ]
                    }
                }
            }
        }
    }
}
`
	template2 := template.Must(template.New("template2.text").Funcs(sprig.TxtFuncMap()).Parse(textTemplate))
	err = template2.ExecuteTemplate(os.Stdout, "template2.text", vals)
	assert.Nil(t, err)
}
