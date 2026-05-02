Apr 28, 2026

## Demo / Walkthrough Roman's Work

Invited [Alex Marple](mailto:alex@sovra.ai) [Alex Maeda](mailto:alexmaeda@sovra.ai) [Jimmy Qian](mailto:jimmy@sovra.ai) [Roman Shteinberg](mailto:roman@sovra.ai) [Patrick Nouvion](mailto:patrick@sovra.ai)

Attachments [Demo / Walkthrough Roman's Work](https://calendar.google.com/calendar/event?eid=NGtnMzEyYTI0MWh2anZ0MTlybnZkaWY0YnMgYWxleEBzb3ZyYS5haQ)

Meeting records [Transcript](https://docs.google.com/document/d/1PFf8w25GqxyDRnZc1-xeSE_4ocBB3luI85TTuayG8xk/edit?usp=drive_web&tab=t.c5eywhpgyg9i) *(Some recordings unavailable)*

### Summary

The team reviewed data pipeline progress and established Subject Data Service ingestion via batch processes.

**Data Pipeline Architecture Review**  
The pipeline integrates standardized canonical representations and the evidence model to ensure robust data quality. Generative AI tests validate edge cases against large datasets.

**Subject Data Service Strategy**  
Attributes will function as first class objects to support multiple independent writers. The team decided that Elastic Search is not the production target for this data.

**Ingestion Mechanism Definition**  
The group decided to utilize batch processes for data ingestion to avoid REST API overhead. Systems will pull parquet files from S3 buckets via asynchronous job endpoints.

### Next steps

- [ ] \[Roman Shteinberg\] Provide Code Pointers: Give navigation pointers to code defining data models. Provide specific examples for the data project.

- [ ] \[Roman Shteinberg\] Change Ingestion Target: Stop implementing Elastic Search ingestion. Determine Subject Data Service storage method and start the new ingestion process toward that target.

- [ ] \[Alex Marple\] Mention Patrick: Post link thread mention to Patrick Nouvion.

### Details

* **Geographical Locations and Time Zones**: The participants began by confirming their geographical locations and corresponding time zones, noting that Roman Shteinberg is six hours ahead of Alex Marple, who is on the West Coast. A clarification was made regarding Washington state, which is on the West Coast, versus Washington D.C., which is on the Eastern Coast. They also briefly discussed the existence of multiple locations named Moscow in the USA ([00:00:00](?tab=t.c5eywhpgyg9i#heading=h.4pkvvzha8zgt)).

* **Meeting Purpose and Scope**: Alex Marple defined the meeting's objective as a 60-minute session to review Roman Shteinberg's progress and work, particularly by walking through what they have been building ([00:00:45](?tab=t.c5eywhpgyg9i#heading=h.kuqvx05ol3kk)). The goal was to provide teammates with a jumping-off point for navigating the code and understanding the past two weeks' efforts, particularly concerning the data project ([00:01:36](?tab=t.c5eywhpgyg9i#heading=h.lh7vz2y8l2gr)). This sharing is intended to align everyone on the path toward the Subject Data Service (SDS) and proactively address potential integration questions related to data models and ingestion methods ([00:02:37](?tab=t.c5eywhpgyg9i#heading=h.3b59zyoj44ue)).

* **Pipeline Structure and Data Flow**: Roman Shteinberg summarized their main contributions, focusing on making the data pipeline solid in terms of data flow, including separate stages for intermediate results. They emphasized the value of having intermediate data representations to facilitate debugging and track decisions made with the data. This approach moves beyond simple backup and normalization to establish a more robust structure ([00:04:27](?tab=t.c5eywhpgyg9i#heading=h.1ukxbeavp6h4)).

* **Data Normalization and Standardization**: A key focus of Roman Shteinberg's work is implementing standardized, canonical representations for data, such as using alpha-2 or alpha-3 standards for naming countries. They also employed libraries for validating telephone numbers (E.164 standard) and parsing addresses using Google-maintained libraries to ensure greater precision and handle edge cases ([00:06:18](?tab=t.c5eywhpgyg9i#heading=h.lgurpmmqgikw)). The strategy is to utilize well-maintained libraries as a baseline for standardization, acknowledging that perfect semantic validation is not fully achieved ([00:08:41](?tab=t.c5eywhpgyg9i#heading=h.4jwa339kboed)).

* **Implementation of the Evidence Model**: Roman Shteinberg introduced the evidence model, which they stated is widely used in data engineering ([00:08:41](?tab=t.c5eywhpgyg9i#heading=h.4jwa339kboed)). The concept of evidence is defined as information derived from the raw data or record, which is distinct from the raw data itself. This distinction allows for better confidence or ranking model development by differentiating data sources, such as a country field from raw data versus a country field derived from a raw address ([00:10:00](?tab=t.c5eywhpgyg9i#heading=h.5w5raa82llaj)).

* **Testing and Alias Management for Canonical Names**: Jimmy Qian inquired about the testing of the new standardization libraries, prompting Roman Shteinberg to explain that they use generative AI to formulate edge cases and write tests based on those cases ([00:11:58](?tab=t.c5eywhpgyg9i#heading=h.648t9mp15twy)). Roman Shteinberg described resolving issues with country naming, particularly non-canonical names like Wales or Scotland, by regulating cases and using the \`pi country\` library alongside an alias table ([00:13:26](?tab=t.c5eywhpgyg9i#heading=h.48031i9j82l2)). The current testing is being conducted on the 50 GB PDL data set, with checks for data distribution planned to identify and remove "garbage" values ([00:15:07](?tab=t.c5eywhpgyg9i#heading=h.m0wd3z75dcu0)).

* **Current Data Pipeline Output**: Alex Marple and Jimmy Qian discussed the output of Roman Shteinberg's pipeline, confirming that it currently outputs parquet files before any Elastic Search ingestion. Roman Shteinberg confirmed that they would inspect those parquet files before the ingestion stage, which is their current focus ([00:17:54](?tab=t.c5eywhpgyg9i#heading=h.12v2n6mvimzr)).

* **Deployment and Execution Environment**: Roman Shteinberg confirmed they are running the data processes on the external server accessed via an account provided by Jimmy Qian ([00:20:02](?tab=t.c5eywhpgyg9i#heading=h.e166z79wg9ol)). They also use a local PC with 32 CPU cores and utilize screen sessions on the server for long-running tasks, such as a 15-hour normalization process ([00:20:52](?tab=t.c5eywhpgyg9i#heading=h.gw9pken6qbwb)). Alex Marple clarified that their questioning about deployment was simply to understand the current operational state, not to impose immediate GitHub deployment ([00:18:52](?tab=t.c5eywhpgyg9i#heading=h.s39ljcdxvoow)).

* **Data Modeling: Subject, Record, and Attribute**: Alex Marple introduced the need to align on the data model for the Subject Data Service, specifically the relationship between Records, Subjects, and Attributes. They reiterated that a Subject is fundamentally an identifier that links to other objects, while Records contain the raw data ([00:23:00](?tab=t.c5eywhpgyg9i#heading=h.xm3pwmxcnoq)). The discussion centered on whether Subjects should have a rigid, field-based structure or if Attributes should be defined as separate, first-class objects attached to a Subject ([00:24:15](?tab=t.c5eywhpgyg9i#heading=h.th1pfw54wao)).

* **Attributes as First-Class Objects**: Alex Marple expressed their opinion that Attributes should be their own first-class object to support multiple, independent writers and allow for selective loading of data ([00:25:13](?tab=t.c5eywhpgyg9i#heading=h.ijhatlx4kqag)) ([00:29:26](?tab=t.c5eywhpgyg9i#heading=h.n9mezwfi4nbk)). This structure means that updating a subject would involve writing separate attribute objects (e.g., an attribute for name and a separate attribute for phone number) ([00:25:13](?tab=t.c5eywhpgyg9i#heading=h.ijhatlx4kqag)) ([00:30:21](?tab=t.c5eywhpgyg9i#heading=h.qgsfubiygpf2)). Jimmy Qian agreed with Alex Marple's analysis and suggested that Alex Marple define the final data model for the team to build toward ([00:27:02](?tab=t.c5eywhpgyg9i#heading=h.sqidvnf0oyru)).

* **Distinction Between Evidence and Attribute**: The group sought clarification on the difference between Evidence and Attributes, especially in the context of Records ([00:39:06](?tab=t.c5eywhpgyg9i#heading=h.kpvrajndzigl)). Roman Shteinberg's Evidence model connects raw data to the normalized values derived from them (e.g., source field, evidence name, and value) ([00:38:15](?tab=t.c5eywhpgyg9i#heading=h.rjrfxbbufqqt)). Attributes, however, are values attached to a Subject, derived from reconciling one or more Records, and may have associated Evidence objects explaining the calculation ([00:40:56](?tab=t.c5eywhpgyg9i#heading=h.urlg95uvz4dl)).

* **Rethinking Elastic Search Ingestion**: A critical decision was reached regarding the role of Elastic Search in the pipeline. Jimmy Qian and Alex Marple both stated that the current Elastic Search environment is for development and investigation, not production, and should not be the authoritative store ([00:55:22](?tab=t.c5eywhpgyg9i#heading=h.x5b499ny4tzz)). Alex Marple concluded that Roman Shteinberg should not worry about ingesting into Elastic Search for now, but should focus on defining the interface for ingesting into the Subject Data Service ([00:56:09](?tab=t.c5eywhpgyg9i#heading=h.17dmd8j2afqt)).

* **Subject Data Service Ingestion Target**: The conversation concluded by defining the Subject Data Service (SDS) as the definitive ingestion target ([00:56:55](?tab=t.c5eywhpgyg9i#heading=h.73o532es9h2j)). Alex Marple clarified that direct writing to the current Postgress database backend of SDS is undesirable; instead, ingestion should be done through the SDS Rest API, potentially via batch upload APIs. The next step is to define the interface between Roman Shteinberg's parquet output and the SDS for ingestion ([00:57:46](?tab=t.c5eywhpgyg9i#heading=h.rmb34yo537c7)).

* **Defining the Data Ingestion Mechanism**: The participants agreed that using a batch process for data ingestion is preferable over repeatedly calling a REST API due to the high overhead associated with the latter approach. Roman Shteinberg initially found the idea of ingesting via HTTP to be a strange approach, which Alex Marple clarified by suggesting the preferred method involves a single call to initiate ingestion, providing a batch ID and a pointer to an S3 file or bucket containing parquet files ([00:58:47](?tab=t.c5eywhpgyg9i#heading=h.ayu3hsa5nim3)). This initiation would then lead to the system pulling the files into the data store ([00:59:33](?tab=t.c5eywhpgyg9i#heading=h.qocxnuty4qfq)).

* **Establishing the Ingestion Workflow**: The proposed workflow involves making a call to a new endpoint with the records and subjects, which would trigger the system to pull the data from S3, meaning the initiation is done by the caller, but the actual data transfer is a pull operation. Alex Maeda summarized the desired process, suggesting that SDS should expose a bulk ingestion job endpoint, receive the parquet location from the pipeline, create an asynchronous job to process it, and track the status, while retaining a smaller, normal REST API for direct, smaller writes. The participants agreed that further discussion is needed to finalize the details of this mechanism ([00:59:33](?tab=t.c5eywhpgyg9i#heading=h.qocxnuty4qfq)).

* **Managing Technical Documentation and Links**: The participants discussed the location of technical documentation, specifically a Mermaid link, which Alex Marple confirmed is sourced in the architecture project and linked on the relevant pull request ([01:00:39](?tab=t.c5eywhpgyg9i#heading=h.8zqp3sovr3d5)). Patrick Nouvion mentioned the difficulty of tracking materials across multiple repositories and asked permission to post links widely in the Slack channel to make them easier to find. Alex Marple agreed to the widespread posting of links, provided that the source of the documentation remains in the architecture repository for cross-cutting items ([01:02:18](?tab=t.c5eywhpgyg9i#heading=h.n604id13jjd4)).

* **Meeting Conclusion**: Roman Shteinberg departed due to a prior appointment, and the meeting concluded with the participants confirming that they had covered the necessary topics. The participants thanked each other and ended the session ([01:02:18](?tab=t.c5eywhpgyg9i#heading=h.n604id13jjd4)).

*You should review Gemini's notes to make sure they're accurate. [Get tips and learn how Gemini takes notes](https://support.google.com/meet/answer/14754931)*

*How is the quality of **these specific notes?** [Take a short survey](https://google.qualtrics.com/jfe/form/SV_9vK3UZEaIQKKE7A?confid=Q4gdMKRIIRyfFpllTCypDxISOAIIigIgABgDCA&detailid=standard&screenshot=false) to let us know your feedback, including how helpful the notes were for your needs.*