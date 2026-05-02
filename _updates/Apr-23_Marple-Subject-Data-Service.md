# Meeting Transcript: Alex Maeda & Alex Marple (Marple)

**Date:** [Insert date]
**Participants:** Alex Maeda, Alex Marple
**Project:** Subject Data Service (SDS)

---

## Opening / Introductions

**Alex Marple:** Alright, cool. Hey Alex, how are you doing? All right, so we've got to decide. Do you have a preferred name? I think we're at least two Alexes now. We might have a third coming on in the future. I don't know if you ever prefer that. Yeah, sorry, go ahead.

**Alex Maeda:** OK.

**Alex Marple:** Yes. Happy to have you guys.

**Alex Maeda:** OK, all right, I can go with Maeda.

**Alex Marple:** A lot of times, when I meet people outside of work, folks are like, "Oh, your name's Alex Marple," and eventually it just becomes Marple as a nickname, which is totally fine with me. So OK, so Marple and Maeda, I think we'll do that.

**Alex Maeda:** Cool. Yeah, so.

---

## Project State & Upcoming Changes

**Alex Marple:** What I was going to say — wanted to say — so two things. One is, yeah, so conscious that things are a bit chaotic. Documentation is sparse and it's also a little bit fluid. So I would just encourage you — feel free to ask questions, ask on Slack, ask the tech channel, ask when you're seeing folks, or reach out. We are conscious that things are a bit chaotic. So yeah.

**Alex Maeda:** OK.

**Alex Marple:** With that said, from the offsite, Jimmy and I — we did have one session where we did a little bit of whiteboarding around how — it's kind of called the search orchestrator — that is probably going to lead to a few changes. I don't think they're large, but they're things that currently don't happen in Subject Data Service right now. Basically one is, I think we're going to need to either add namespaces to subjects so that you could have like, oh, these are the subjects written by the internal data service versus subjects that we're pulling from Pitbull, and you could select only one of them or search within only one of them or pull from only one of them. So I think that's one thing. Not super difficult in terms of adding it. It's not there right now and we haven't written this down yet. So we want to try to get a little bit more up to date and written and shared tomorrow or by tomorrow.

And the second thing is — we are very likely going to have duplicate subjects, at least for right now. Because the subject orchestrator may need to write the same subject multiple times if they get searched on Monday versus Friday versus the following Monday, and we get new data. So both of those things are likely going to happen. Since we both agree that there's a little bit of a smell to creating duplicate data just because of the way that we're kind of like extracting data on the fly, while we're doing it, we would have to decide between creating just multiple records in the same table or if we want to separate the table where an internal process is writing subjects and they're written consistently or written once and then not modified very often from the table where it's going to grow relatively quickly because you have a lot of similar subjects that get written and maybe they have different IDs but they're all kind of related to the same person. So I suspect that those are going to be at least different tables, but they don't have to be. So the two of us can think on what we want to do there.

**Alex Maeda:** Yeah.

**Alex Marple:** I realize that action is on me to get that written so that it's a little bit more clearly described — like why we're adding this. Because I still don't have a better proposal yet. I still agree that it's kind of a weird structure, but I think it's the best thing that we have right now. But before getting to that, I want to give you some time to see — I don't know if you have any questions, comments, questions about the task, about the existing codebase. Because I'm conscious I stood this up last week. Have not put very much thought into the service — there's very little functionality in there. I've put a little bit more thought into the data model — like what are records, what are subjects, why are they structured the way that they are. So that has a little bit more. But all of this I consider version 0, so we can make updates to it. So that was a bit of a word salad as an intro. Do you have any things you want to talk about or any particular questions or things that will help you better understand what's going on in the codebase or how to work with Roman or Patrick or any of that kind of stuff?

---

## Alex Marple's BAM Walkthrough

**Alex Marple:** Yeah, OK. Let me maybe take 15 minutes to walk through a little bit of the stuff on the BAM side, like the behavioral analysis module. I'm going to throw a lot of information at you. You don't need to retain all of this. I will come back at the end and kind of describe them a little bit more. I'm going to say bluntly here. Let me do a couple things.

**Alex Maeda:** Perfect. Yeah, walk me through it. I'll take notes mentally, but I appreciate you circling back at the end. I find that helps me separate important to remember from good context to have. And don't worry about information overload. I'd rather you give me too much than too little. I can always ask clarifying questions. Ready when you are. Fire away.

**Alex Marple:** Let me first share my screen. Anyway, that's fine. OK, so you can do that. And again, apologies, I'm going to throw a bunch of links at you. So OK, so if you — I assume you're on the Sovra Drive. If not, within the shared Google Drive for the team — do you have access to this? I assume yes, but if not I will work to get you added.

**Alex Maeda:** OK, cool.

**Alex Marple:** OK, cool. So within — I'm not going to go through all the structure, but within this 07 teams folder, there's basically a folder for each team. UX engineering is combined but then we get split out into our own thing. This is not all of the documents that are related to engineering, but for the ones that are across engineering that are in Drive — we had a conversation eight weeks ago and no one can remember exactly what was said and maybe we can find the transcript. You'll see that. As an example, this subject search and detail page requirements — this is actually not the one that I actually want to go into. I'm just noting that the conversations that we had in the offsite going through the MVP user flow on Figma. I'm going to jump away from this. It's just an example of the sheet that I was using to try to capture the things that we were discussing with the user flows around the search and detail page as something that's a little bit more structured. This is one where I think we're doing an OK job of it.

But if you look within that folder, there's also this called ESPAM July 1 beta, which I'm putting things for. So if you look at that, ESPAM pipelines — it's a draw.io file. That's what I'm going to reference.

**Alex Maeda:** OK.

---

### ESPAM Pipeline & CE Features

**Alex Marple:** OK. So if you've pulled this up — this is an unfortunately large file. I'm not going to go through all of this, but the idea here is that our science team, Dr. Bluestein, has basically described ESPAM, which is this pipeline that we can do to extract different components of behavioral analysis from these subjects. This relatively long pipeline — this diagram is attempting to capture that in a little bit more detail than is currently in the patents. The patents are not a design document, but they do describe the process that Doc is thinking through. This diagram is an attempt to be a little bit more specific about which of those boxes are describing an artifact versus a process. For each process, what are its inputs, what are its outputs, what's the order in which all these things happen, where are humans involved? So I don't necessarily need to go through the whole thing, but if you're able to see this, you'll see that the box at the top says preprocessing. 160 different traits in about 70 categories — things that we can assess about a subject, behavioral or psychological features that we can assess about the subject.

If my screen share is still working, I'm showing this other sheet — this rating template VR. There's a fair amount that's in here, but the sheet "normal personality 44" contains 44 different features in this normal personality category and describes what someone who rates low — yeah, OK. So this is basically that line that you're looking at, particularly within normal personality, there is this trait called N2 anger and hostility. This has a trait definition. There is a low, medium, and high benchmark for this trait. So someone who rates low for anger and hostility is even-tempered, slow to anger, forgiving, maintains composure when frustrated, gives others benefit of the doubt... medium and high, so on and so forth. So each of these is one of those CE features.

The first step of that pipeline is basically looking to produce for a particular subject — what is the rating on a scale of 1 to 5 where 1 is low, 3 is medium, 5 is high for that particular feature for that subject? As well as reasoning for giving that rating and the evidence that we want to cite for giving that rating. Evidence is basically a particular social media post, a list of social media posts, or this online blog, or this article written about that person. This evidence is basically the record. So if we're coming back to SDS, this is the record that supports that reasoning or the list of records that supports that reasoning. And then a confidence score which is either 1, 2, or 3 — or Low, Medium, or High.

I'll note that confidence technically only has Low, Medium, and High, whereas the rating has 2 and 4 — you actually can say this is between low and medium and that would be a score of 2, or this is between medium and high as a score of 4. Confidence doesn't have that, but for the sake of trying to keep the scores being the same, I basically say both of these are on a scale of 1 to 5 and confidence just happens to be either 1, 3, or 5 — you can't have a 2 or 4. But anyway, when we're talking about ratings, we're basically talking about the output of that first step of the BAM pipeline, which is for this subject — like for Alex Marple, what is my rating and the reasoning behind that rating and the confidence for that rating and evidence supporting that rating for each of those 160 traits. That's what the feature ratings are.

If we keep going from top to bottom through that pipeline — the BAM pipeline — there are other steps. So after assessing ratings, there's these components where we say, OK, well, based on ratings, let's identify what are the risks and the strengths related to this candidate. If you are looking at that rating sheet, you'll also see that there are tabs for risk taxonomy and strength taxonomy. So these are the categories of risks that we believe subjects can have, the categories of strengths that subjects can have. So we would try to assign — similar to this — a hypothetical future state where we say, OK, well, maybe we have an LLM-based judge that is assigning or extracting risks and strengths, as well as a traditional ML like XGBoost-based pipeline that's trying to do the same thing, and then we need to reconcile the scores.

And the further and further you get down this pipeline, the more abstract the definition of exactly what the input and output and the internals of each of these components are. Doc has definitely put thought into this, but we understand the top part pretty well and we know that we want to build pieces of that by July. The further down you go, the later in time this is probably going to happen unless it's something that we need to build for the MVP. So you'll see that all over this are these — I'm basically using color coding to mean blue means I have some idea, red means we have no idea what this is — it needs further description. And sometimes I'll put those red notes just meaning, hey, there's open questions here. We do need to drill down on these when we start to build them. But I at least understand the intention of what Doc is looking to do for these. You don't need to absorb all of this, but the reason that I wanted to share it is that that's what the ratings are. Those CE ratings stands for Conceptualization Engine — the terminology Doc gave. The first step of this pipeline is assigning those ratings for subjects. So that's what ratings are.

If you're looking at the BAM prototype codebase — the one that I've been working on a bit more — the piece of the Subject Data Service model is based on that. You'll see that that does include things like salient scores and confidence and all these things. I've — it's clear that it's just not well defined yet.

Let's go to inflection points, maybe do risk. Yeah, OK. Here's an example — I'm not sure if you can see this. Here's an example — I've defined this class "Pattern of Life." The only thing in the docstring is "TODO" — we don't know what this is. We need to find the schema for this. So there's plenty of structures like Pattern of Life or MitigatorAmplifierResult or even — maybe not RiskScalingScore but RiskMeaning. We have a little bit more specifics that a RiskMeaning has a what, a why, a how, and a when. I have no idea how Doc intends to capture it. I don't know if this is free-form text or if it's pulling from a taxonomy of subjects that could be added. I don't know what this looks like. We do know that there are those five components. This will change as we get there. I just wanted to share this as an example that this is version 0.X of this data schema. Don't treat anything here as super, super stable. Some of this is just representing — we know that this is coming, it's an object that exists. I don't want to drop that on the floor, but since it's not critical to what we're doing right now, let's not waste time trying to over-specify something that we don't even have the inputs needed to try to compute again. So I hope that answers your question about CE features as well as ratings. But if it doesn't, feel free to tell me what else would be useful to try to answer that question or would help you put some of this stuff into context.

---

### BAM vs. SDS Split

**Alex Marple:** The reason that a little bit of the stuff sits in BAM but was moved over to SDS — the short answer is that the BAM prototype was basically a project where I was attempting to put together a prototype of everything that was being shouted out or mentioned as something that we might do in an implementation so that we could start to make some progress, get some feedback, see what we would need to do to build something like that. The piece that I think stays within BAM is all the stuff related to the BAM pipeline — the temporal workflows that model the process better that are described in the patents. Those I think are going to stay in BAM, use LLMs to accomplish those, we can use traditional ML, we can do things that are deterministic. But I think all of that stays in BAM.

However, some of the other things that are modeled here are things that are needed by the chatbot but aren't necessarily part of that BAM pipeline. The best example of that is subjects and the raw data that we are using to chat about subjects. Keep going back to — where does Adam Bunkadeco live? Or has he ever run for public office? The BAM pipeline isn't deciding that — that's a fact that you can just go to the records for. So I needed an API to be able to say, hey, create subjects, associate records with subjects, search over the records for subjects. And I was just throwing everything into a monorepo because discussions were kind of all over the place. But now that we actually have more of a plan for what the backend architecture is going to look like, that information about what subjects exist and what information do we have about subjects and being able to either look up records or potentially search records related to a subject — that we think belongs in its own service, and that service is the Subject Data Service. So I intend to get most of that moved over.

---

### Subjects & Records Definition

**Alex Marple:** I think that the two entities that exist — subjects and records — for the most part we've decided on the terminology. We're saying that the things that Sovra deals with are subjects right now. Those are people. We're using subjects because in the future those might be LLCs, or a partnership, or a household, or something maybe other than people. But right now, basically those are people. People are subjects.

Records for us are the external data that we are pulling in. Records basically are everything from social media profiles to social media posts to court records to marriage records to information that we're pulling off of breach data sets that Grant is going and pulling. The idea is that a record represents that raw data. There's a couple components to that. Let me go over to Subject Data. So I'm looking at the most recent version of this. Oh, that's bad — I still call it trace data. OK, yeah. So for this you can see — record has metadata. So it has the subject ID that it's related to — this is basically a secondary index. There's metadata — information related to how the record was obtained. So this is like, oh, we fetched it from this URL, we fetched it at this date, we're tagging it as the run from April whatever — information on how we pulled it.

The data source — let me confirm that I'm using the right terminology. Yeah, OK. So the data source is our attempt to describe what this data is using enums rather than making it completely unstructured. So platform is a list of places where we know that we're pulling data, and content type is the classification of the content type. So we might say our content types are social profiles, social posts, criminal records, marriage records, court records, news articles. I forget where the list of that is in the Go implementation, but the idea is that there's a fixed set of content types that we define.

And then the platform data basically has structure to say — OK, there's going to be some raw data. The raw data is whatever we got for that record. This is unstructured. It can vary from platform to platform, it can vary from record to record. But the idea is that there's also this normalized field that says — OK, based on the content type, we should be able to extract particular properties. An example would be — for a social media profile, if that's a content type, we should be able to try to extract the profile picture and the profile description as text from a social media profile regardless of what platform it came from and the structure that the raw data has. So the idea is that each of these records, even though it has this unstructured raw data, has this normalized field that basically says — we want to be able to extract the same information for records that have the same content type.

You also see that I put these images and videos as a placeholder here. At some point we probably will start doing some image and video analysis, basically saying — for each of the images related to this record, can we extract a description and the subjects in it and the setting and the apparent activity and the presentation style? We did a little bit of prototyping with basically throwing this at some vision models and trying to extract this. We haven't really used this very much anywhere. So I do anticipate that what is in this media analysis will change. But the idea is that if records do contain non-text data — like if they contain media — we would want to do some processing and extract stuff. We're not doing any of that yet, but that's something that you could attach to a record.

So subjects being a person that we can find know things about, and records — those are going to continue to exist and they're going to keep their name in Subject Data Service.

---

### Attributes / Properties Definition

**Alex Marple:** The third object that isn't in the API yet, but that I think we want to add — and this is going back to what I was originally describing — is what I used to call properties. I think we've realigned on calling them attributes, but both are super generic names. I'll try to just describe what those are. A property or an attribute is information that we attach to a subject. We say — OK, subject has a list of properties or a list of attributes.

An attribute has a name. The name of a property might be "current address." It has a type. The type of a current address — maybe we say it's a string, maybe we say this is an address that has street line one, street line two, optional apartment number, stuff like that. Based on that type, it would also contain the structured data that adheres to that type. So you might say — a subject has a current address property, but it can also have a current employer property, or an employment history property, or a marriage status property, or a relationship status property. We haven't defined exactly what those properties are. As an example — is relationship status a part of marriage status, or are they two separate properties? We don't have a good taxonomy of exactly what those properties are. But the idea is that properties are information that we can attach to a subject.

My mental model is that the reason that we're not attaching them to subjects themselves is that there isn't a single source of truth for what is writing those properties. I'll give an example. We've talked about things like FFIP — what Tim Jones had mentioned — where we say, OK, if I want to determine someone's marriage status, these are the things I checked. I want to be able to fetch these types of records and then based on the contents of these records, I can assign marriage status being likely married with a confidence of 60%, whatever. That is one thing that we can do to assign properties to subjects. But we could also look at the BAM ratings as properties that we're attaching to subjects. That's a separate pipeline that goes and does this offline analysis, and it's a potentially long and involved process. And then at the output, it says — I'm going to attach these properties back to a subject. So the idea is that we do want to be able to look up all of the properties for a subject, but those properties can be assigned by multiple writers. And based on the type of the property, data that's in there is going to have a different schema, and it can be assigned a confidence interval that doesn't exist in the API yet or doesn't exist in the data model yet. But that's the idea of one of the other concepts that needs to get added to Subject Data Service.

---

### Trace Data Rename Note

**Alex Marple:** And I guess the one thing that I wanted to note — this isn't super relevant, but my OCD is getting to me — is that the reason that this is named trace_data.go is that originally in the BAM prototype I named the object TraceData. TraceData is the raw data that we get from external data sources. But when we came back together, Roman and Jimmy had been using the term "record," so I said — OK, that's fine, it's just a name. So I renamed them to record. However, I renamed all the instances of trace data here rather than the file name, so that's why this is still called trace_data.go. This really should just be called record.go or records.go. That's not super important, but just a little bit of history to why this is still called trace_data.go — that should be changed.

---

### Use Cases: Patrick (Frontend) & Roman (Identity Resolution)

**Alex Marple:** But I'll include this thing because it's good to have for the end of the transcription. There's at least two use cases for SDS right now, or at least for the read APIs right now. The read use case is — Patrick is working on the frontend as well as this search orchestrator thing. One of the things that he's going to need to do is when he's building out the detail cards on search results — so you say, I search for John Smith and you've got John Smith and Jonathan Smith and John Smith Senior — for each of those he's going to say, OK, this is a subject, I need to fetch the information to render in that card. So that would be their location and their name and their age and the obfuscated phone number and obfuscated email — all that kind of stuff. So the intention is that card can be rendered by reading from Subject Data Service and saying — here's a subject, and I want to read this list of properties. I don't have an opinion right now on whether you provide a list of properties to get only those back, or we just say — here's all the properties for the subject. I don't have an opinion on the right way to start, but that's the first read use case for SDS.

But I will note that there's also this use case for the chatbot where the chatbot might say — if someone is asking "is this subject married?" It probably wants a tool that says — let me go to SDS, let me look up the marriage status property for the subject, and go and check to see what the contents of that is. So there'll be multiple callers reading from SDS. But Patrick's probably the one to talk to.

Roman — what Roman is working on is on the right side. He is basically going through all of the breach data and the data that Grant has been scraping to try to say — OK, I have a record that this email is attached to this physical address, and this physical address is also attached to this username, and this username is attached to this first name/last name pair and phone number — to basically try to aggregate that into saying, oh hey, based on these 20 different records, we actually believe this is John Smith. Here's some of the information about him. So that's doing two things. It's saying — hey, I've identified a subject, I have enough information to indicate that this is the same person, so I want to create a subject. That would be writing the subject. And it's also saying — OK, I believe I've found a first name, a last name, an address, a phone number — all these other properties or attributes or whatever we end up calling them — and either write those as properties, or write a record that says — hey, here's my analysis, this intermediate state where I think all these are true, let me write that record so that some analysis engine can go in and either just copy them over and trust it or reconcile it with other information. Not sure exactly what that looks like yet.

I mentioned that because I'm not sure exactly what the output is going to look like of Roman's process. I don't know if he's running an ETL job, or he's got a Pandas notebook that's going to spit out a giant list of subjects and records, or if he's doing it incrementally and there's a service at runtime that's going to spit them out one by one. I don't know. But we'll need to work with him just to confirm that we can get the things that he's producing ingested into SDS in a way that makes sense for other consumers to be able to read it.

I say I'm not sure exactly what's going on there because part of this will be confirming that the service and the API that we're building matches up with what Roman needs. I wouldn't be surprised if it turns out — hey, we actually need a bulk import. You call an operation saying — here's an S3 file, I want you to ingest it — and we say — cool, here's a job ID, I'll update the job ID when I'm done. Maybe we do something like that. I don't know that it's needed, I'm just conscious that that might be something we do based on the conversation with Roman.

**Alex Maeda:** Yeah.

---

## Closing & How to Reach Alex Marple

**Alex Marple:** All right, cool. Any additional questions? I think you said you were good. We've got stuff in the transcript. I just wanted to kind of close back around on why to talk to Patrick and Roman as well.

**Alex Maeda:** OK, all right. Cool.

**Alex Marple:** Cool. All right. Yeah. Again, if you have any questions, I'm on Slack. I'm generally on like 9 to 6 — that's my general hours. Sometimes I might be on late or off early, on earlier, or off later, but generally 9 to 6 I'm available. And also I'll say Slack is probably the best way to reach me. And I will say — never worry about sending me a Slack message at the wrong time. I manage my notifications. If I'm away then I won't pay attention to it, but never worry about sending me a message. I will never be annoyed that I got a message regardless of the time of day. I can't guarantee I'll respond immediately, but I will never be annoyed that you sent a message.

**Alex Maeda:** Thanks, Marple. This was really helpful — seriously.

---

**End of Transcript**