# Data Objects

There is a collaborative copy of this document at [https://docs.google.com/document/d/1DO8sMHEGS3RGQbPqY1wuOPscxKE3OlY\_dDscbbWe30o/edit?usp=sharing](https://docs.google.com/document/d/1DO8sMHEGS3RGQbPqY1wuOPscxKE3OlY_dDscbbWe30o/edit?usp=sharing)  
There is a version-controlled copy of this document at [https://github.com/sovraai/bam\_prototype/blob/pipeline-planning/docs/data-model.md](https://github.com/sovraai/bam_prototype/blob/pipeline-planning/docs/data-model.md)

This doc is a work in progress. As of 2026-03-04, this only contains data objects relevant for [ML Classification](https://drive.google.com/file/d/1KvC6m5kInm48id-V6qcSoUXfK16Kd0eU/view) and [Risk Enumeration](https://drive.google.com/file/d/1tWikU1QWrbWbFZLSIj5QyEDljNiGcgHv/view). [Contextualization](https://drive.google.com/file/d/1OGRHRa7i7UOKYPFBx06EzR8cXX5OrDob/view) is coming next.

## CE Feature Definition

* feature: FeatureId|string \- the feature  
* category: CategoryID|string \- the category in which the feature is grouped: Normal Personality, Personality Under Stress, Character, Psychopathy, Values, Motivations, Aptitude, (maybe Risk Taxonomy)  
* subcategory: SubcategoryId|string \- a subcategory within the category. For Normal Personality, this is either Big Five or the HEXACO dimension. For Risk Taxonomy, this is "Main Risk Category".  
* definition: string \- description of the feature  
* low\_benchmark: string \- description of a subject with a rating of Low (1)  
* medium\_benchmark: string \- description of a subject with a rating of Medium (3)  
* high\_benchmark: string \- description of a subject with a rating of High (5)

**Example:**

```json
{
   "feature": "Courage",
   "category": "Character",
   "definition": "Sustained perseverance under prolonged adversity; quantified via endurance and refusal to disengage.",
   "behavioral_indicators_low": "Withdrawal or disengagement as adversity persists; energy collapses over time.",
   "behavioral_indicators_medium": "Intermittent persistence with visible strain; relies on external encouragement.",
   "behavioral_indicators_high": "Long-horizon persistence with adaptive pacing; adversity integrated rather than avoided."
}
```

## Subject CE Feature

* feature: FeatureId|string \- the feature  
* category: CategoryID|string \- the category in which the feature is grouped: Normal Personality, Personality Under Stress, Character, Psychopathy, Values, Motivations, Aptitude, (maybe Risk Taxonomy)  
* subcategory: SubcategoryId|string \- a subcategory within the category. For Normal Personality, this is either Big Five or the HEXACO dimension. For Risk Taxonomy, this is "Main Risk Category".  
* rating: string \- rating as a string, L, L/M, M, M/H, H  
* rating\_as\_number: int \- rating 1-5  
* rating\_reasoning: string \- why the rating was assigned  
* evidence: list- identifiers for the evidence to support the reasoning. This should be the ID for a single trace data source.  
* confidence: int \- confidence as a string, L, M, H  
* confidence\_as\_number: int \- confidence 1/3/5

**Example:**

```json
{
   "subject_id": "subj_123456",
   "feature": "Courage",
   "category": "Character",
   "rating": "Moderate to high",
   "rating_as_number": 4,
   "rating_reasoning": "continually perseveres towards goals, even in the face of failure. Two failed congressional bids, now pursuing state office",
   "evidence": [
      "a2",
      "b30"
   ],
   "confidence": "Medium",
   "confidence_as_number": 3,
   "confidence_reasoning": "very public well documented failures and continued pursuits"
}
```

## Subject Risk Feature

Fields inherited from Subject CE Features

* probability: float \- the probability this risk is relevant for the subject

**TODO:** is score still required? Is it meaningful?

**Example:**

```json
{
   "subject_id": "subj_123456",
   "feature": "Courage",
   "category": "Character",
   "probability": 0.7,
   "rating": "Moderate to high",
   "rating_as_number": 4,
   "rating_reasoning": "continually perseveres towards goals, even in the face of failure. Two failed congressional bids, now pursuing state office",
   "evidence": [
      "a2",
      "b30"
   ],
   "confidence": "Medium",
   "confidence_as_number": 3,
   "confidence_reasoning": "very public well documented failures and continued pursuits"
}
```

## Interaction Feature Definition

* interaction\_label: string \- name for the interaction  
* feature: FeatureId|string \- the interacting feature  
* hexaco\_modulators: list- the HEXACO feature(s) that modulate  
  * feature: FeatureId|string \- the feature  
  * condition: the condition for the rating of the feature. TBD, how granular do these conditions need to be?  
* personality\_under\_stress\_modulator: list\<FeatureId|string\> \- the Personality Under Stress feature(s) that modulate  
* context\_modulators: string \- the subject context that modulates. TODO, is this free-text, or is there a list of possible values?  
* psychopathy\_modulators: list\<FeatureId|string\> \- the Psychopathy feature(s) that modulate  
* character\_modulators: list\<FeatureId|string\> \- the Character feature(s) that modulate  
* interaction: InteractionType|string \- Moderator, Attenuator, Configurational effect, Effect/Outcome, Interaction, Three-way profile  
* interaction\_description: string \- description of the interaction  
* applied\_value: string \- description of the outcome  
* priority: int \- priority 1-5  
* strength: int \- strength 1-5  
* confidence: int \- confidence 1/3/5  
* notes: string \- free form notes, used to cite related research, mention other concepts  
* level: LevelType|string \- level at which interaction occurs. TBD, how is this used?  
* direction: string \- direction, TBD, how is this used?  
* outcome: RiskTaxonomy|string \- the risk from the taxonomy for that can result from this interaction. TBD some of the outcome values in the sheet are free text, not taxonomy risks  
* model\_encoding: string \- implementation suggestion for this interaction, presumably for the ML \+ Explainability path

**Example:**

```json
{
   "interaction_label": "Self-Regulation Buffer: Substance Use",
   "feature": "NEO:N",
   "hexaco_modulators": [
      {
         "feature": "HEX:Conscientiousness",
         "condition": "x == 5"
      }
   ],
   "personality_under_stress_modulators": [],
   "context_modulators": ["Addictive history"],
   "psychopathy_modulators": [],
   "character_modulators": [
      {
         "feature": "CHAR:Hope (protective)",
         "condition": null
      }
   ],
   "interaction": "Moderator",
   "interaction_description": "High HEX:Conscientiousness tends to reduce the link between high N and problematic substance use by adding structure, planning, and inhibition of emotion-driven impulses.",
   "applied_value": "Substance-use risk; Relapse considered low in these circumstances",
   "priority_1_5": 2,
   "strength": 1,
   "confidence_across_settings": 3,
   "notes": "Turiano et al. (2012; MIDUS; N≈4000): interaction β ≈ −.28 to −.35; higher C attenuates N→substance use. Use 'attenuates' rather than 'eliminates'.",
   "level": "Trait×Trait",
   "direction": "Attenuating interaction",
   "outcome_risk_taxonomy": "Psychological Risks",
   "model_encoding": "Linear interaction"
}
```

## Subject Interaction Feature

**TODO:** Is there any subject-specific information for interaction features?

## Subject Areas of Risk

* string \- A single free-text field describing the risks for the subject

**Example:**

Areas of Risk:  We don't have a lot of data in general, but we don't have information to suggest that he is able to integrate and work well with a team.  He has no recent involvement with law enforcement, which would indicate he is likely stable from a behavioral standpoint, not prone to impulsivity or taking significant risks.

## Subject Inflection Point

## Subject Demographics and Biopsychsocial History

* string \- A single free-text field describing the candidate. Demographics include age, gender, education, socioeconomic status, cultural background. Biopsychosocial includes physical, mental, or social context.

**TODO:** what does "good" look like for demographics and biopsychosocial?

**Example:**

Candidate is a 63 y/o married male with 1 child.  He and his wife currently reside in the United States.  He was born and raised within a wealthy family in the suburbs of Paris.  Mr. Combes attended prestigious institutions and has had a rather successful career within the telecommunications and financial industries.  He has frequently been placed in key leadership situations and is sought after to be a part of organizational boards.  His wife, Christie Julien is a well known concert pianist.  She is active on social media, where she has at times embraced and promoted right-wing and conspiratorial ideas.  She does not reference her husband or his work in any of her postings.  He has a very limited social media presence and only rarely posts / references family.

## Subject Conclusion/Synthesis

* string \- A single free-text field summarizing strengths, weaknesses, and risks for the subject

**Example:**

Limited information on this candidate.  We have no relational or occupational history.  Limited data would indicate there are not significant risks.  He appears to somewhat independent and motivated to achieve / succeed through his entrepreneurial efforts.

## Risk Classification Weights

* w\_LLM: float \- the weight to scale LLM risk probabilities by  
* w\_ML: float \- the weight to scale ML risk probabilities by

**Example:**

```json
{
  "w_LLM": 0.7,
  "w_ML": 0.2
}
```

## Risk Classification Disagreement Threshold

* threshold: float \- the threshold for difference of absolute value of risk probabilities, above which we don't consider consensus reached for the risk

**Example:**

```json
{
  "threshold": 0.15
}
```

## Subject Classified Risk Profile Output

* feature: FeatureId|string \- the feature  
* category: CategoryID|string \- the category in which the feature is grouped: Normal Personality, Personality Under Stress, Character, Psychopathy, Values, Motivations, Aptitude, (maybe Risk Taxonomy)  
* subcategory: SubcategoryId|string \- a subcategory within the category. For Normal Personality, this is either Big Five or the HEXACO dimension. For Risk Taxonomy, this is "Main Risk Category".  
* probability: float \- the fused probability (based on ML probability and LLM probability)  
* calibrated\_confidence: float \- confidence as a float (computed using Platt scaling)  
* consensus: boolean \- whether the probabilities from the ML and LLM path differed by more than the disagreement threshold  
* shap\_feature\_attribution: string \- a summary of the ML-assigned risk probability  
* llm\_rating\_reasoning: string \- why the rating was assigned by LLM  
* llm\_evidence: list- identifiers for the evidence to support the reasoning. This should be the ID for a single trace data source.

**Note:** this object is similar to "Subject Risk Feature", but is defined separately because they serve different purposes.

**Example:**

```json
{
   "subject_id": "subj_123456",
   "feature": "Courage",
   "category": "Character",
   "probability": 0.7,
   "calibrated_confidence": 0.732,
   "consensus": true,
   "llm_rating_reasoning": "continually perseveres towards goals, even in the face of failure. Two failed congressional bids, now pursuing state office",
   "llm_evidence": [
      "a2",
      "b30"
   ],
  "shap_feature_attribution": "TBD, what is a good example of this?"
}
```

## Profile-derived Risk

**TODO:** I don't know what this is or what its schema looks like. I suspect that this is no different from "Subject Risk Feature", or possibly "Subject Classified Risk Profile Output"

## Categorized Risk

**TODO:** I suspect that this is the same schema as "Subject Classified Risk Profile Output"

## Habits

**TODO:** I don't know what this looks like

## Lifestyle Signals

**TODO:** I don't know what this looks like

## Self-reported Context

**TODO:** I don't know what this looks like

## Environment & Situational Factors

**TODO:** I don't know what this looks like

## Pattern-of-Life Output

**TODO:** I don't know what this looks like

## Emergent Risk

**TODO:** I suspect that this is the same schema as "Categorized Risk" (which I suspect is in turn the same as "Subject Classified Risk Profile Output")

## Risk Salience Score

* salience\_score: float \- a measure of how salient (prominent among all risks) the risk is for the subject

**Example:**

```json
{
  "salience_score": 0.15
}
```

## Risk Lifestyle Relevance

* salience\_score: float \- a measure of how salient (prominent among all risks) the risk is for the subject

**Example:**

```json
{
  "lifestyle_relevance": 0.23
}
```

## Risk Confidence Rating

* Confidence: float \- a measure of confidence in the risk for the subject

**Example:**

```json
{
  "confidence": 0.72
}
```

## Risk Meaning

* what: string \- plain language description. TBD, does this have structure?  
* why: string \- contextual to subject's life or user's lifestyle. TBD, does this have structure?  
* how: string \- behavioral manifestations. TBD, does this have structure?  
* when: date \- date of last event. TBD, is this more than just a single date?  
* tension: list- a list of the shadow linkages between risks and strengths  
  * strength: str|dict \- the strength (or an ID by which to find the strength)  
  * risk: str|dict \- the risk (or an ID by which to find the risk)

**Example:**

```json
{
  "what": "TBD, need an example",
  "why": "TBD, need an example",
  "how": "TBD, need an example",
  "when": "2026-01-28T13:30:00.000Z",
  "tension": {
    "linkages": [
      {
        "strength": "TBD: is this a strength feature ID, a strength ID, or a complete strength object?",
        "risk": "TBD: is this a risk feature ID, a risk ID, or a complete risk object?"
      }
    ]
  }
}
```

**TODO:** I don't know what this looks like

## Subject Scored Risk Output

* subject\_classified\_risk\_profile\_output: Subject Classified Risk Profile Output \- the risk before scoring  
* composite\_score: float \- the composite score  
* salience\_score: Risk Salience Score \- salience  
* lifestyle\_relevance: Risk Lifestyle Relevance \- lifestyle relevance  
* confidence\_rating: Risk Confidence Rating \- confidence  
* meaning: Risk Meaning \- meaning  
* provenance: dict \-  TBD whether this is any different from the reasoning and evidence in “Subject Classified Risk Profile Output”

**Note:** this object is similar to "Subject Risk Feature", but is defined separately because they serve different purposes.

**Example:**

```json
{
  "subject_classified_risk_profile_output": {
    "subject_id": "subj_123456",
    "feature": "Courage",
    "category": "Character",
    "probability": 0.7,
    "calibrated_confidence": 0.732,
    "consensus": true,
    "llm_rating_reasoning": "continually perseveres towards goals, even in the face of failure. Two failed congressional bids, now pursuing state office",
    "llm_evidence": [
      "a2",
      "b30"
    ],
    "shap_feature_attribution": "TBD, what is a good example of this?"
  },
  "composite_score": 0.82,
  "salience_score": {
    "salience_score": 0.15
  },
  "lifestyle_relevance": {
    "lifestyle_relevance": 0.23
  },
  "confidence_rating": {
    "confidence": 0.72
  },
  "meaning": {
    "what": "TBD, need an example",
    "why": "TBD, need an example",
    "how": "TBD, need an example",
    "when": "2026-01-28T13:30:00.000Z",
    "tension": {
      "linkages": [
        {
          "strength": "TBD: is this a strength feature ID, a strength ID, or a complete strength object?",
          "risk": "TBD: is this a risk feature ID, a risk ID, or a complete risk object?"
        }
      ]
    }
  }
}
```

## Life Event Definition

* event\_type: LifeEventType \- enum values: JobChange, RelationshipChange, Financial, Legal, ...  
* definition: string \- what types of events constitute this change?

## 

**TODO:** flesh this out further (or confirm this is sufficient)

**Example:**

```json
{
  "event_type": "JobChange",
  "definition": "A change in employment, either the start of a job or the end of one"
}
```

## Subject Life Event

* event\_type: LifeEventType \- enum values: JobChange, RelationshipChange, Financial, Legal, ...  
* reasoning: string \- an explanation of why the life event was detected  
* evidence: list- identifiers for the evidence to support the reasoning. This should be the ID for a single trace data source.

**TODO:** flesh this out further (or confirm this is sufficient)

**Example:**

```json
{
  "event_type": "JobChange",
  "reasoning": "The subject left their employer after 9 years",
  "evidence": [
    "linkedin_1234"
  ]
}
```

## Resilience Measurement

* behavioral\_response: dict \- response pattern pre- & post-event. TBD how this is structured  
  * before: string \- description of response before event  
  * after: string \- description of response after event  
* correlated\_activity: string \- which behavioral patterns changed in temporal proximity to the event  
* recovery\_trajectory: string \- how and how quickly behavioral patterns return to pre-event baseline levels  
* new baseline: string \- did the person return to the same baseline or establish a new one  
* ce\_prediction\_alignment: dict \- TBD how this is structured  
  * ce\_change: dict \- does the subjects CE profile after the event differ from before? TBD how this is structured  
  * action\_change: string \- do the subjects actions after the event differ from their CE profile

**TODO:** What are behavioral patterns? **TODO:** How are behavioral patterns measured? How are they baselined?

**Example:**

```json
{
  "behavioral_response": {
    "before": "TBD, need an example",
    "after": "TBD, need an example"
  },
  "correlated_activity": "TBD, need an example",
  "recovery_trajectory": "TBD, need an example",
  "new_baseline": "TBD, need an example",
  "ce_prediction_alignment": {
    "ce_change": {
      "TBD": "need structure and example"
    },
    "action_change": "TBD, need an example"
  }
}
```

## Inflection Point Property Tags

* direction: Direction|string \- risk-accelerating or protective  
* magnitude: float \- how much the trajectory actually shifted)  
* permanence: Permanence|string \- initial classification as likely permanent or likely temporary, subject to dynamic reclassification  
* ce\_features: CE Feature Definition \- which CE features and domains are affected  
* confidence: float \- how certain the system is that this is a genuine inflection point versus noise

**Example:**

```json
{
  "direction": "protective",
  "magnitude": 0.3,
  "permanence": "likely temporary",
  "ce_features": [
    {
      "feature": "Courage",
      "category": "Character"
    }
  ],
  "confidence": 0.85
}
```

## Subject Trace Data

* metadata: information about how the trace was obtained  
  * website: string \- website where the trace was obtained (this looks redundant with data\_source.platform)  
  * url: string \- full URL where the trace was obtained  
  * capture\_date: date \- when the trace was obtained  
  * tags: list- free-form text tags  
  * identifiers: list- additional identifiers we can assign  
* data\_source: information on what type of trace this is  
  * platform: Platform|string \- the source of the trace, e.g. "instagram", "twitter", "website". Supports a fixed set of values  
  * content\_type: ContentType|string \- a classification of this type of content, shared across platforms, e.g. "post", "profile", "comment"  
* platform\_data: the content of the trace, split into platform-specific data and "normalized" data in a format shared by traces of the same data\_source.content\_type  
  * normalized: map \- normalized schema for the content\_type. TBD on these schemas  
  * raw: dict \- map \- schema specific to the platform \+ content\_type. Possibly schema-less  
* images: a list of images  
  * image  
    * url: string \- where to access the image content. Note that this is different from the public URL, which can change or be removed entirely.  
    * context: ContextType|string \- how the image was used. Tentative values: "default" by default, "original" to indicate the original image (i.e. the one being reacted to), "reaction" to indicate the reaction  
  * image\_analysis: the information to extract from the image (TBD, the fields below are suggestions based on lora\_demo)  
    * description: string \- a description of the image  
    * subjects: string \- people in the image  
    * setting: string \- the location of the subjects in the image  
    * apparent\_activity: string \- activities in which the subjects are engaged in the image  
    * presentation\_style: string \- how the image is presented (TBD)  
    * mood\_tone: string \- mood of the image (TBD)  
    * text\_in\_image: string \- text features in the image (intended to capture text on a reposted flyer, but not the street signs in a photo)  
* videos: a list of videos  
  * video  
    * url: string \- where to access the image content. Note that this is different from the public URL, which can change or be removed entirely.  
    * context: ContextType|string \- how the image was used. Tentative values: "default" by default, "original" to indicate the original image (i.e. the one being reacted to), "reaction" to indicate the reaction  
  * video\_analysis: the information to extract from the video (TBD, the fields below are just a copy of the fields for image)  
    * description: string \- a description of the video  
    * subjects: string \- people in the video  
    * setting: string \- the location of the subjects in the video  
    * apparent\_activity: string \- activities in which the subjects are engaged in the video  
    * presentation\_style: string \- how the video is presented (TBD)  
    * mood\_tone: string \- mood of the video (TBD)  
    * text\_in\_video: string \- text features in the video (intended to capture text on a reposted flyer, but not the street signs in a photo)

**Data Source × Data Type Matrix (for March/April/May study):**

| Data Source | Access | Time Window | Profile Text | Profile Image/Analysis | Post Text | Post Image/Analysis | Post Video/Analysis | Reaction Text | Reaction Image/Analysis | Reaction Video/Analysis | Reaction Original Text | Reaction Original Image/Analysis | Reaction Original Video/Analysis |
| :---- | :---- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| Bluesky | public | As many as available | Y | Y | Y | Y | N | N | N | N | N | N | N |
| Crunchbase | public | As many as available | Y | Y | n/a | n/a | n/a | n/a | n/a | n/a | n/a | n/a | n/a |
| Facebook | public | As many as available | Y | Y | Y | Y | N | N | N | N | N | N | N |
| FlightRadar24 | public | As many as available | n/a | n/a | ? | n/a | n/a | n/a | n/a | n/a | n/a | n/a | n/a |
| Glassdoor | public | As many as available | n/a | n/a | Y | n/a | n/a | n/a | n/a | n/a | n/a | n/a | n/a |
| Google | public | As many as available | n/a | n/a | ? | Y | N | n/a | n/a | n/a | n/a | n/a | n/a |
| Instagram | public | As many as available | Y | Y | Y | Y | N | N | N | N | N | N | N |
| LinkedIn | public | As many as available | Y | Y | Y | Y | N | N | N | N | N | N | N |
| Reddit | public | As many as available | Y | N | Y | Y | N | N | N | N | N | N | N |
| Threads | public | As many as available | Y | Y | Y | Y | N | N | N | N | N | N | N |
| TikTok | public | As many as available | Y | Y | Y | Y | N | N | N | N | N | N | N |
| Twitter | public | As many as available | Y | Y | Y | Y | N | N | N | N | N | N | N |
| Bankruptcy Records | public |  | n/a | n/a | Y | ? | n/a | n/a | n/a | n/a | n/a | n/a | n/a |
| Civil Court Records | public |  | n/a | n/a | Y | ? | n/a | n/a | n/a | n/a | n/a | n/a | n/a |
| Criminal Court Records | public |  | n/a | n/a | Y | ? | n/a | n/a | n/a | n/a | n/a | n/a | n/a |
| Family Court Records | public |  | n/a | n/a | Y | ? | n/a | n/a | n/a | n/a | n/a | n/a | n/a |
| Mortgage/Deed Records | public |  | n/a | n/a | Y | ? | n/a | n/a | n/a | n/a | n/a | n/a | n/a |

Y \= yes, N \= no, n/a \= column not meaningful for this source, ? \= undecided

**Note on "Reaction" columns**: For social media, "reaction" refers to the subject's engagement with others' content (retweets, shares, replies, comments). "Reaction Original" is the original content that the subject reacted to — providing context for *what* they chose to engage with. For court/public records, reactions are not applicable.

**Note on court/public records**: "Post Text" maps to the document text content (filings, rulings, docket entries). "Post Image/Analysis" maps to scanned documents or exhibits. Video is generally not available from public record sources.

**Note**: discussion in Science / Tech Sync 2026-03-09 for scope of March/April/May study

* No reactions  
* Pull as many posts as we can if cost isn’t an object.  
  * Observe the cost of the pull for reference  
* No medical records  
* No video

**Example:**

```json
{
   "metadata": {
      "website": "instagram.com",
      "url": "https://www.instagram.com/p/DTWBVT_EilZ",
      "capture_date": "2026-01-10T00:00:00",
      "tags": [],
      "identifiers": []
   },
   "data_source": {
      "platform": "instagram",
      "content_type": "post"
   },
   "platform_data": {
      "raw": {
         "page_type": "post",
         "username": "adembunkeddeko",
         "post_date_absolute": "January 10, 2026",
         "post_date_relative": "1w",
         "likes_count": null,
         "caption_text": "We're building a people-powered campaign for New York State Comptroller \u2014 and every contribution counts. With our first fundraising deadline at 11:59pm tomorrow, now is the moment to chip in. Donate via the link in our bio.",
         "caption_mentions": [],
         "caption_hashtags": []
      },
      "normalized": {
         "text": "We're building a people-powered campaign for New York State Comptroller \u2014 and every contribution counts. With our first fundraising deadline at 11:59pm tomorrow, now is the moment to chip in. Donate via the link in our bio.",
         "reaction_count": 0
      }
   },
   "images": [
     {
       "image": {
          "url": "https://hosted-copy-of-data/image_abcdef.jpg",
          "context": "default|original|reaction"
       },
       "image_analysis": {
         "description": "Instagram post featuring an individual outdoors holding a sign",
         "subjects": "A person wearing a hoodie and jacket, standing outdoors near trees and a snowy path, holding a sign that reads 'HOW FUNDRAISING WORKS FOR THIS COMPTROLLER RACE'",
         "setting": "Outdoor park or similar area during winter, with bare trees and snow on the ground",
         "apparent_activity": "The person is participating in a fundraising campaign for a comptroller race, encouraging contributions via a link in their bio",
         "presentation_style": "Single photo post with accompanying text caption",
         "mood_tone": "Informative and motivational, aiming to engage viewers in a political campaign",
         "text_in_image": "adembunkeddeko We're building a people-powered campaign for New York State Comptroller \u2014 and every contribution counts. With our first fundraising deadline at 11:59pm tomorrow, now is the moment to chip in. Donate via the link in our bio."
       }
     }
   ],
   "video": [
      {
         "video": {
            "url": "https://hosted-copy-of-data/video_123456.mp4",
            "context": "default|original|reaction"
         },
         "video_analysis": {
            "description": "Instagram post featuring an individual outdoors holding a sign",
            "subjects": "A person wearing a hoodie and jacket, standing outdoors near trees and a snowy path, holding a sign that reads 'HOW FUNDRAISING WORKS FOR THIS COMPTROLLER RACE'",
            "setting": "Outdoor park or similar area during winter, with bare trees and snow on the ground",
            "apparent_activity": "The person is participating in a fundraising campaign for a comptroller race, encouraging contributions via a link in their bio",
            "presentation_style": "Single photo post with accompanying text caption",
            "mood_tone": "Informative and motivational, aiming to engage viewers in a political campaign",
            "text_in_image": "adembunkeddeko We're building a people-powered campaign for New York State Comptroller \u2014 and every contribution counts. With our first fundraising deadline at 11:59pm tomorrow, now is the moment to chip in. Donate via the link in our bio."
         }
      }
   ]
}
```

**TODO:** Defining the schema for all Data Sources and types of content needs a deep dive. This will define things like platform\_data, image\_analysis, video\_analysis, and context.

**TODO:** For data sources with large amounts of text (court records, deeds), we'll probably want to extract a summary.

## S-VIP

**TODO:** Is there any existing material to draw from on the data format for S-VIP?

## S-PACE

**TODO:** Is there any existing material to draw from on the data format for S-PACE?