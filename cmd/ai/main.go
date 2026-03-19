package main

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langgraphgo/graph"
)

// 1. 定义状态
type AgentState struct {
	Messages []string
}

func main() {
	// 2. 创建图
	g := graph.NewMessageGraph()

	// 3. 添加节点 (Node)
	g.AddNode("agent", func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		// 模型思考并决定是否调用工具
		// state.Messages = append(state.Messages, "AI: 我建议搜索天气...")
		state = append(state, llms.TextParts(llms.ChatMessageTypeAI, "我建议搜索天气..."))
		return state, nil
	})

	// 4. 设置入口点和边
	g.SetEntryPoint("agent")
	g.AddEdge("agent", graph.END) // 暂时简单结束

	// 5. 编译并运行
	runnable, _ := g.Compile()
	ms, err := runnable.Invoke(context.Background(), []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeHuman, "text")})
	fmt.Println(ms, err)
	// runnable.Invoke(context.Background(), AgentState{})

	chat, err := ollama.New(ollama.WithCustomTemplate("deepseek-r1:latest"))
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := chat.Call(context.Background(), "你是谁？")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
