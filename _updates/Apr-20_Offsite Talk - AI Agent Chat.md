# 2026-04-20 \- Offsite Talk \- AI Agent Chat

## Prep Checklist

1. Confirm we’re recording

## 1a. What is an "agentic" chatbot? (\~5 min)

**Key message**: An agentic chatbot doesn't just answer — it *researches*, then answers.

* A **traditional chatbot** takes a question, sends it to an LLM, and returns whatever the LLM says in one shot.  
    
  - It can only use what it already "knows" from training data.  
  - If it doesn't know something, it guesses (hallucination).


* An **agentic chatbot** has a loop:  
    
  1. Receive user question  
  2. Decide: "Can I answer this, or do I need to look something up?"  
  3. If it needs data, call a **tool** (e.g., look up a subject's risk profile)  
  4. Read the tool's result  
  5. Decide again: "Do I have enough to answer, or do I need more?"  
  6. Repeat until ready, then respond


* **Analogy**: A traditional chatbot is like asking someone a question from memory. An agentic chatbot is like asking a researcher — they'll go look things up, cross-reference sources, and then give you a grounded answer.  
    
* The agent decides **what to do next** at each step. This is what makes it "agentic" — it has agency over its own process.

---

## 1b. The Levers We Control (\~15 min)

**Framing**: "These are the things we tune to make the chatbot better. Think of them as dials on a mixing board."

### Lever 1: The Model (\~3 min)

* The LLM (Large Language Model) is the "brain" of the agent.  
* Different models have different capabilities and costs:  
  - **Larger models** (70B+ parameters): Better at reasoning, following complex instructions, and using tools correctly. More expensive, slower.  
  - **Smaller models** (8B parameters): Faster, cheaper, but may misunderstand instructions or pick the wrong tool.  
* BAM supports multiple providers: Together AI, vLLM, Ollama, RunPod.  
* Model choice affects:  
  - Quality of the final answer  
  - Whether it picks the right tools  
  - Cost per conversation  
  - Response speed  
* **Key point**: This lever has the **least differentiation** — everyone in the market has access to the same models. Our advantage comes from the next three levers.

### Lever 2: The System Prompt — Instructions (\~4 min)

* The system prompt is the **"job description"** we give the agent at the start of every conversation.  
* It tells the agent:  
  - Who it is: *"You are a behavioral assessment assistant"*  
  - What tools it has and when to use them  
  - How to behave: *"Always ground your answers in the data returned by tools. Cite specific traits, scores, or evidence."*  
  - What to do first: *"Use get\_subject\_summary as a starting point when exploring a new subject."*  
* **This is one place where we encode domain expertise.** For example:  
  - "When the user asks about safety, prioritize the risk profile"  
  - "When citing trait ratings, explain what the score means"  
  - "Always include caveats about the limitations of the data"  
* Changing the prompt is **fast and free** — no retraining, no code changes, immediate effect.  
* The prompt is an important lever because it embeds domain knowledge that is relevant for every question (use cases we support, application of S-BAM material, site navigation, etc.)

### Lever 3: The Tools — What the Agent Can Do (\~5 min)

* Tools are the agent's "hands" — functions it can call to retrieve or compute information.  
    
* Each tool has three parts:  
    
  1. **Name**: What the LLM sees (e.g., `get_ce_features`)  
  2. **Description**: A plain-English explanation of what the tool does — this is what the LLM reads to decide *when* to use it  
  3. **Schema**: What inputs the tool accepts (e.g., `subject_id`, `category`)


* The BAM chatbot currently has **20 tools** in different states of (a) implementation and (b) having data available. They are currently organized into four groups:

| Tool | Theme | Description | Key Inputs | Implemented | Data |
| :---- | :---- | :---- | :---- | :---: | :---: |
| `list_subjects` | Subject | List all subjects with IDs and display names | — | Yes | Yes |
| `lookup_subject` | Subject | Subject metadata: name, trace data count, ratings count, platforms | `subject_id` | Yes | Yes |
| `search_trace_data` | Trace Data | Semantic search over a subject's trace data | `subject_id`, `query`, `max_results?` | Yes | Yes |
| `get_ce_features` | CE Ratings | LLM-rated CE features for a subject (1-5 scale), optionally filtered by category | `subject_id`, `category?` | Yes | Yes\* |
| `get_risk_profile` | Risk | Classified risk profiles with consensus and calibrated confidence | `subject_id` | Yes | No |
| `list_ce_feature_definitions` | Taxonomy | List CE feature definitions (the rating rubric), optionally filtered by category | `category?` | Yes | Yes |
| `get_ce_feature_definition` | Taxonomy | Full definition of one CE feature including 1-5 rating benchmarks | `feature_name`, `definition_id?` | Yes | Yes |
| `list_risk_definitions` | Taxonomy | List risk taxonomy definitions (subset of CE features in the "Risk Taxonomy" category) | `subcategory?` | Yes | No |
| `list_life_event_definitions` | Taxonomy | Life event type definitions used for inflection point detection | — | No | No |
| `list_interaction_definitions` | Taxonomy | Trait interaction definitions (how traits modulate each other) | `feature?`, `interaction_type?` | No | No |
| `get_scored_risks` | Risk Deep-Dive | Scored risks with composite score, salience, lifestyle relevance, and meaning narrative | `subject_id` | Yes | Yes |
| `get_risk_mitigators_amplifiers` | Risk Deep-Dive | Amplifying and mitigating factors for a specific risk | `risk_key` | No | No |
| `get_pattern_of_life` | Risk Deep-Dive | Behavioral patterns and lifestyle signals for a risk | `risk_key` | No | No |
| `get_inflection_points` | Inflection Points | Life events for a subject, optionally filtered by event type | `subject_id`, `event_type?` | No | No |
| `get_ce_feature_detail` | CE Detail | Full evidence chain for one feature: reasoning, quotes, confidence | `subject_id`, `feature_name` | Yes | Yes |
| `compare_subject_features` | CE Detail | Side-by-side comparison of multiple CE features for a subject | `subject_id`, `feature_names` | No | No |
| `list_trace_data` | Trace Data | Browse trace data entries by platform/content type | `subject_id`, `platform?`, `content_type?`, `limit?` | No | Yes |
| `get_trace_data_entry` | Trace Data | Full content of a specific trace data entry by ID | `trace_id` | Yes | Yes |
| `list_categories` | Overview | The 9 assessment categories with trait counts and names | — | No | Yes |
| `get_subject_summary` | Overview | Dashboard of data availability across all dimensions for a subject | `subject_id` | Yes | Yes |

* **Quality of the tool description** directly affects whether the agent picks the right tool. A vague description \= wrong tool chosen.  
    
* **Quality of the tool's output** directly affects the quality of the answer. If a tool returns messy or incomplete data, the answer suffers.  
    
* **Adding a new tool \= teaching the agent a new skill.** For example, we could add a "compare\_subjects" tool to enable side-by-side analysis.

### Lever 4: Playbooks / Routing (\~3 min)

* Beyond individual tools, we can create **different agents** (different prompt \+ tool combinations) for different tasks.  
* BAM already does this:  
  - **Agent Chat**: Full behavioral assessment assistant with 12 tools  
  - **Subject Search View**: Translates natural language into search API queries — completely different prompt, no assessment tools  
* Same underlying LLM, but different instructions and capabilities depending on the task.  
* **Future direction**:  
  - A "router" agent that automatically picks the right specialist based on the user's question  
  - Embedding "playbooks" in the system prompt — e.g., "for safety questions, always follow these steps: 1\) check risk profile, 2\) check psychopathy traits, 3\) check trace data for concerning content"  
  - Sub-agents that handle specific categories of questions with specialized expertise

---

## 1c. How It All Fits Together — The Loop (\~5 min)

**Walk through a concrete example**: "Where does Adem live?"

Below is an example of the messages sent to/from the Chatbot LLM by the end of the conversation:

1. The system prompt, which is the first message given to the LLM. You can see it in `llm_messages[0]`, or at [https://github.com/sovraai/bam\_prototype/blob/main/backend/services/prompt/prompt\_builder.py\#L575-L601](https://github.com/sovraai/bam_prototype/blob/main/backend/services/prompt/prompt_builder.py#L575-L601), or extracted below:  
   1. You are a behavioral assessment assistant with access to tools that let you look up subject data, social media trace data, CE feature ratings, risk profiles, scored risks, and the CE/risk taxonomy definitions.  
   2.   
   3. When the user asks a question, decide whether you can answer from the conversation context alone or whether you need to call one or more tools first. If you need data, call the appropriate tool. After receiving tool results you may call another tool or respond to the user.  
   4.   
   5. You have access to both subject-specific tools (ratings, risks, trace data) and taxonomy tools (CE feature definitions with rating benchmarks, risk definitions). Use the taxonomy tools to explain what a rating means or to list available features. Use get\_subject\_summary as a starting point when exploring a new subject.  
   6.   
   7. Always ground your answers in the data returned by tools. Cite specific traits, scores, or evidence when possible.  
   8.   
   9. IMPORTANT: When you have finished calling tools and have the data you need, you MUST respond to the user with a substantive answer that synthesizes the tool results. Never respond with an empty message.  
   10.   
   11. IMPORTANT: All tools that accept a subject\_id parameter require the exact subject ID string as returned by list\_subjects. Do NOT guess or abbreviate subject IDs — always use the exact ID from list\_subjects.  
2. The user’s first message (“Where does Adem live?”). You can see it in `llm_messages[1]`.  
3. The agent’s request to invoke a tool (search\_trace\_data). You can see it in `llm_messages[2]`.  
4. The result of invoking the tool, which is provided back to the agent. In this case, it is a list of TraceData similar to the query “Adem location”. You can see it in `llm_messages[3]`.  
5. Finally, there is the response produced by the chatbot for the user. You can see this in `action` & `message`.

```
[
  {
    "llm_cost_audit": {
      "workflow_id": "",
      "entries": [
        {
          "id": "19ed0405-0ff0-42fe-8ca7-ca420eb919b8",
          "workflow_id": "",
          "activity_name": "agent_plan_next_action",
          "provider": "together",
          "model": "meta-llama/Llama-3.3-70B-Instruct-Turbo",
          "input_tokens": 2597,
          "output_tokens": 42,
          "estimated_cost_usd": 0.00232232
        }
      ],
      "total_input_tokens": 2597,
      "total_output_tokens": 42,
      "total_cost_usd": 0.00232232
    },
    "action": "respond",
    "message": "Adem lives in Brooklyn, NY. This information was found on his Instagram profile and various news articles and campaign websites, which mention his campaign for New York State Comptroller and his presence in Brooklyn.",
    "tool_name": null,
    "tool_args": {},
    "llm_messages": [
      {
        "role": "system",
        "content": "You are a behavioral assessment assistant with access to tools that let you look up subject data, social media trace data, CE feature ratings, risk profiles, scored risks, and the CE/risk taxonomy definitions.\n\nWhen the user asks a question, decide whether you can answer from the conversation context alone or whether you need to call one or more tools first. If you need data, call the appropriate tool. After receiving tool results you may call another tool or respond to the user.\n\nYou have access to both subject-specific tools (ratings, risks, trace data) and taxonomy tools (CE feature definitions with rating benchmarks, risk definitions). Use the taxonomy tools to explain what a rating means or to list available features. Use get_subject_summary as a starting point when exploring a new subject.\n\nAlways ground your answers in the data returned by tools. Cite specific traits, scores, or evidence when possible.\n\nIMPORTANT: When you have finished calling tools and have the data you need, you MUST respond to the user with a substantive answer that synthesizes the tool results. Never respond with an empty message.\n\nIMPORTANT: All tools that accept a subject_id parameter require the exact subject ID string as returned by list_subjects. Do NOT guess or abbreviate subject IDs — always use the exact ID from list_subjects.\n\nThe current subject being discussed has subject_id: \"843400dd-bcdc-42ea-ae46-db838712c287\". Use this ID when calling tools unless the user asks about a different subject."
      },
      {
        "role": "user",
        "content": "Where does Adem live?"
      },
      {
        "role": "assistant",
        "content": null,
        "tool_calls": [
          {
            "id": "call_1",
            "type": "function",
            "function": {
              "name": "search_trace_data",
              "arguments": "{\"subject_id\": \"843400dd-bcdc-42ea-ae46-db838712c287\", \"query\": \"Adem location\", \"max_results\": 10}"
            }
          }
        ]
      },
      {
        "role": "tool",
        "tool_call_id": "call_1",
        "content": "{\"results\": [{\"trace_id\": \"3259e937-da8d-416f-a3bd-6079b49d694c\", \"platform\": \"website\", \"content_type\": \"post\", \"content\": \"website post: A campaign landing page for Adem Bunkeddeko, a candidate for New York State Comptroller. Campaign promotion and engagement Professional and informative ADEM BUNKEDDEKO FOR NEW YORK STATE COMPTROLLER Democrat for NYS Comptroller. Check out our website! Donate to the campaign! Show your support with our social media toolkit! Follow us on Medium! Press NY1 City & State New York Join ademfornewyork on Linktree Tags: source:B - AB2.pdf:32\", \"url\": \"https://linktr.ee/ademfornewyork\", \"tags\": [\"source:B - AB2.pdf:32\"], \"relevance_score\": 0.318}, {\"trace_id\": \"d6389807-b51c-41ec-83e7-fbf0c05be260\", \"platform\": \"website\", \"content_type\": \"post\", \"content\": \"website post: A campaign landing page for Adem Bunkeddeko, a candidate for New York State Comptroller. Campaign promotion and engagement Professional and informative ADEM BUNKEDDEKO FOR NEW YORK STATE COMPTROLLER Democrat for NYS Comptroller. Check out our website! Donate to the campaign! Show your support with our social media toolkit! Follow us on Medium! Press NY1 City & State New York Join ademfornewyork on Linktree Tags: source:C - AB3.pdf:99\", \"url\": \"https://linktr.ee/ademfornewyork\", \"tags\": [\"source:C - AB3.pdf:99\"], \"relevance_score\": 0.315}, {\"trace_id\": \"cf10c078-3bce-46a8-bdb8-d677f879e912\", \"platform\": \"instagram\", \"content_type\": \"post\", \"content\": \"instagram post: Tomorrow, December 19th, join us at @cafe.erzulie for the Eclipse Holiday Party! Link to tickets in the bio! Instagram post featuring an event poster for 'Eclipse Holiday Drive' hosted by Cafe Erzulie. Promotion of a holiday-themed music event scheduled for December 19th. Festive and energetic, emphasizing the holiday spirit and live music experience. ECLIPSE HOLIDAY DRIVE HOSTED BY DAVIDFOREIGNFINESSE FINDRSKEEPERS V. TINO CAFE ERZULIE 894 BROADWAY, BROOKLYN, NY 11206 FEATURING \", \"url\": \"https://www.instagram.com/p/DSbF2atjzV7\", \"tags\": [\"Business\", \"Warning\", \"Illustration\", \"Vectors\", \"Text\", \"Retro\", \"Symbol\", \"Sign\", \"Print\", \"Bill\\nVintage\", \"No\", \"Person\", \"Rubberize\", \"Desktop\", \"Danger\", \"Offense\", \"Caveat\", \"Market\", \"Paper\\nLabel\", \"source:C - AB3.pdf:68\"], \"relevance_score\": 0.271}, {\"trace_id\": \"b2c99ef7-0c05-4e4a-bb1c-098420f356e1\", \"platform\": \"news_article\", \"content_type\": \"document\", \"content\": \"news_article document: Adem Bunkeddeko talks bid for state comptroller A screenshot of a Spectrum News article featuring a video interview with Adem Bunkeddeko, a Democratic state comptroller candidate. Adem Bunkeddeko is participating in an interview on 'Inside City Hall,' discussing his bid for state comptroller. Professional and informative, focusing on political news and campaign discussions. Democratic state comptroller candidate Adem Bunkeddeko joined 'Inside City Hall.' (Spectrum News NY1\", \"url\": \"https://ny1.com/nyc/all-boroughs/inside-city-hall/2025/12/05/adem-bunkeddeko-talks-bid-for-state-comptroller\", \"tags\": [\"source:B - AB2.pdf:56\"], \"relevance_score\": 0.271}, {\"trace_id\": \"3a6536e1-00fd-4eea-b67f-1749afa05901\", \"platform\": \"news_article\", \"content_type\": \"document\", \"content\": \"news_article document: Adem Bunkeddeko talks bid for state comptroller A screenshot of a Spectrum News article featuring a video interview with Adem Bunkeddeko, a Democratic state comptroller candidate. The candidate is participating in an interview segment titled 'Inside City Hall,' where he discusses his campaign plans and goals. The tone is informative and professional, typical of news coverage focusing on political campaigns. Democratic state comptroller candidate Adem Bunkeddeko joined 'Ins\", \"url\": \"https://ny1.com/nyc/all-boroughs/inside-city-hall/2025/12/05/adem-bunkeddeko-talks-bid-for-state-comptroller\", \"tags\": [\"source:C - AB3.pdf:123\"], \"relevance_score\": 0.267}, {\"trace_id\": \"5ee2d3b6-0f00-45e6-b99c-c1b4d3c59b82\", \"platform\": \"website\", \"content_type\": \"post\", \"content\": \"website post: A campaign donation page for Adem T. Bunkeddeko, candidate for New York State Comptroller. Campaign fundraising effort Serious and motivational Adem T. Bunkeddeko for New York State Comptroller\\nDonate to Adem's Campaign!\\nNew York isn't just facing an affordability crisis \\u2014 it's facing an opportunity crisis.\\n[...]\\nChoose an amount: $50, $25, $10, $5\\nMake it monthly!\\nCheckout\\nContribution rules\\nPlatform paid for by ActBlue (actblue.com) Tags: source:C - AB3.pdf:158\", \"url\": \"https://secure.actblue.com/donate/smalldollarwknd\", \"tags\": [\"source:C - AB3.pdf:158\"], \"relevance_score\": 0.254}, {\"trace_id\": \"7019ce1d-b48b-4f37-b99f-d5688c939695\", \"platform\": \"instagram\", \"content_type\": \"profile\", \"content\": \"instagram profile: Democrat for NYS Comptroller.\\nhttps://linktr.ee/ademfornewyork Instagram profile page showing a grid of posts related to political campaigning and personal moments. Displaying a collection of posts that highlight the user's political campaign for NYS Comptroller and personal engagements. Professional yet approachable, emphasizing community engagement and personal connection. adem bunkeddeko, 30 posts, 2.1k followers, 908 following\\nDemocrat for NYS Comptroller.\\nhttps://linktr.e\", \"url\": \"https://www.instagram.com/adembunkeddeko\", \"tags\": [\"source:C - AB3.pdf:31\"], \"relevance_score\": 0.238}, {\"trace_id\": \"b140f581-3f53-4af3-82fe-3397767c15f8\", \"platform\": \"website\", \"content_type\": \"post\", \"content\": \"website post: A campaign webpage for Adem T. Bunkeddeko running for New York State Comptroller. Political campaign promotion and informational content about the candidate and the comptroller position. Serious and informative, aiming to convey trustworthiness and expertise. ADEM T. BUNKEDDEKO FOR NEW YORK STATE COMPTROLLER, Democrat for State Comptroller, The Office of the New York State Comptroller, New York isn't just facing an affordability crisis. It is facing an opportunity crisis. Tags: sou\", \"url\": \"https://www.ademfornewyork.com/\", \"tags\": [\"source:B - AB2.pdf:78\"], \"relevance_score\": 0.237}, {\"trace_id\": \"1999b858-8ba8-409b-949b-821d4c4d5acd\", \"platform\": \"instagram\", \"content_type\": \"profile\", \"content\": \"instagram profile: Democrat for NYS Comptroller.\\nhttps://linktr.ee/ademfornewyork Instagram profile page showing a grid of posts related to political campaigning and personal moments. Displaying a mix of political campaign activities and personal moments through Instagram posts. Professional and motivational with a personal touch, emphasizing political engagement and community involvement. adem bunkeddeko, 30 posts, 2.1k followers, 908 following\\nDemocrat for NYS Comptroller.\\nhttps://linktr.ee/ad\", \"url\": \"https://www.instagram.com/adembunkeddeko\", \"tags\": [\"source:B - AB2.pdf:26\"], \"relevance_score\": 0.225}, {\"trace_id\": \"00b6eb95-7a0d-43c1-936e-fab4f2928538\", \"platform\": \"website\", \"content_type\": \"post\", \"content\": \"website post: A campaign donation page for Adem T. Bunkeddeko, candidate for New York State Comptroller. Encouraging donations to a political campaign Serious and motivational Adem T. Bunkeddeko for New York State Comptroller\\nDonate to Adem's Campaign!\\nNew York isn't just facing an affordability crisis \\u2014 it's facing an opportunity crisis.\\n[...]\\nChoose an amount: $50, $25, $10, $5\\nMake it monthly! Checkout\\n[...] Tags: source:B - AB2.pdf:91\", \"url\": \"https://secure.actblue.com/donate/smalldollarwknd\", \"tags\": [\"source:B - AB2.pdf:91\"], \"relevance_score\": 0.224}], \"total_searched\": 140}"
      }
    ],
    "model_override": null
  }
]
```

**Key points to emphasize**:

* The agent made **1 tool call** before answering a single question  
* The tool call was a **deliberate choice** — the agent decided what information it needed  
* The final answer is **grounded in data**, not hallucinated  
* Every LLM call and tool call has a measurable **cost** (tokens in, tokens out, time)  
* The conversation runs as a **Temporal workflow** — it's durable. If the user disconnects and comes back, the conversation continues right where it left off. It also gives us visibility into the agent’s usage

## More Chat Demos

1. Agent Chat: “Has Adem run for public office?”  
   1. Factual lookup based on records. The tool is based on vector search, and the agent answers well.  
2. Agent Chat: “Would Adem enjoy dancing?”  
   1. The agent doesn’t know how the tools available can be used to answer this question.  
3. Subject Search View: “I'm looking for John in Houston”  
   1. This is a different agent, built to assist the user as they search for subjects.  
   2. The stack that supports the chat is the same; the only things that changed are (a) the system prompt, and (b) the tools or lack thereof.

---

## Listing Chat Use Cases

1. [https://docs.google.com/spreadsheets/d/1MwX0SFYGedobb0zAz3AJJXdjvOMU0U5D1ED0UcrdHSE/edit?usp=sharing](https://docs.google.com/spreadsheets/d/1MwX0SFYGedobb0zAz3AJJXdjvOMU0U5D1ED0UcrdHSE/edit?usp=sharing)  
   1. Sheet: “Chat Use Cases”  
      1. Discuss list of cases  
      2. Prioritize the list of cases  
      3. Discuss the evaluation criteria that we construct collectively (Engineering, UX, Science, Marketing)  
   2. Sheet: “Chatbot Questions \- Behavior (WIP)”  
      1. Time permitting, this is the start of evaluation criteria for the Behavioral questions.

