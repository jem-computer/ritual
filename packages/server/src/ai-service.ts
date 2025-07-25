// ABOUTME: AI service for executing prompts using AI SDK
// ABOUTME: Supports multiple models and providers

import { generateText } from 'ai';
import { openai } from '@ai-sdk/openai';
import { anthropic } from '@ai-sdk/anthropic';
import { z } from 'zod';

// Configuration schema
export const AIConfigSchema = z.object({
  openaiApiKey: z.string().optional(),
  anthropicApiKey: z.string().optional(),
  defaultModel: z.string().default('claude-3-5-sonnet-20241022'),
  maxTokens: z.number().default(1000),
  temperature: z.number().default(0.7),
});

export type AIConfig = z.infer<typeof AIConfigSchema>;

// Model mapping
const modelProviderMap: Record<string, string> = {
  // Anthropic models
  'claude-3-5-sonnet-20241022': 'anthropic',
  'claude-3-5-haiku-20241022': 'anthropic',
  'claude-3-opus-20240229': 'anthropic',
  'claude-3-sonnet-20240229': 'anthropic',
  'claude-3-haiku-20240307': 'anthropic',
  // OpenAI models
  'gpt-3.5-turbo': 'openai',
  'gpt-4': 'openai',
  'gpt-4-turbo': 'openai',
  'gpt-4o': 'openai',
  'gpt-4o-mini': 'openai',
};

export class AIService {
  private config: AIConfig;

  constructor(config: Partial<AIConfig> = {}) {
    this.config = AIConfigSchema.parse({
      openaiApiKey: process.env.OPENAI_API_KEY,
      anthropicApiKey: process.env.ANTHROPIC_API_KEY,
      ...config,
    });
  }

  async executePrompt(prompt: string, model?: string): Promise<{
    output: string;
    model: string;
    usage?: {
      promptTokens: number;
      completionTokens: number;
      totalTokens: number;
    };
  }> {
    const modelToUse = model || this.config.defaultModel;
    const provider = modelProviderMap[modelToUse];

    if (!provider) {
      throw new Error(`Unsupported model: ${modelToUse}`);
    }

    try {
      let result;
      
      if (provider === 'anthropic') {
        if (!this.config.anthropicApiKey) {
          throw new Error('Anthropic API key not configured');
        }
        
        result = await generateText({
          model: anthropic(modelToUse),
          prompt,
          maxTokens: this.config.maxTokens,
          temperature: this.config.temperature,
        });
      } else if (provider === 'openai') {
        if (!this.config.openaiApiKey) {
          throw new Error('OpenAI API key not configured');
        }

        result = await generateText({
          model: openai(modelToUse),
          prompt,
          maxTokens: this.config.maxTokens,
          temperature: this.config.temperature,
        });
      } else {
        throw new Error(`Provider ${provider} not implemented`);
      }

      return {
        output: result.text,
        model: modelToUse,
        usage: result.usage ? {
          promptTokens: result.usage.promptTokens,
          completionTokens: result.usage.completionTokens,
          totalTokens: result.usage.totalTokens,
        } : undefined,
      };
    } catch (error) {
      console.error('AI execution error:', error);
      throw new Error(`Failed to execute prompt: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  async testConnection(): Promise<boolean> {
    try {
      const result = await this.executePrompt('Say "hello" in one word.');
      return result.output.toLowerCase().includes('hello');
    } catch (error) {
      console.error('AI connection test failed:', error);
      return false;
    }
  }
}

// Singleton instance
let aiService: AIService | null = null;

export function getAIService(): AIService {
  if (!aiService) {
    aiService = new AIService();
  }
  return aiService;
}

export function initAIService(config?: Partial<AIConfig>): AIService {
  aiService = new AIService(config);
  return aiService;
}