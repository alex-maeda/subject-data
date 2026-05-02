Regarding the "Subject Data Service" ...
The current project is https://github.com/sovraai/subject-data
- I would start with building and running this locally
- Next I would check that you're able to log in to the "sovra org" AWS account (770826159442), view the AWS Organization to see the AWS account for subject-data-beta, and then switch roles into that account. Jimmy or I will need to create a user for you, it shouldn't be a blocker to running locally.
I think the first task is adding "attributes" to the API, which includes:
Updating the schema to include attributes, which are attached to subjects
- Attributes have a type, and a schema, which is determined by the type
- Attribute data may be stored or passed around as an envelope, and adherence to the schema is checked elsewhere
 * Schemas could be versioned to support change
 * Within a version, maybe we can do add-only changes
- We want to incorporate "confidence" or "confidence ranges".
 * Is this applied at the attribute level?
 * Or is this applied to nested fields within the attribute?


I would start with getting that project pulled, run locally, and looking at the API (/docs should load the Swagger UI). Once you have that, we can chat a bit more about the task of adding "attributes". The idea is that "attributes" are the "verification" information produced by S-VIP (if you've seen/heard S-VIP mentioned). The task is to design the API changes to add attributes to SDS.

For people to talk to related to this:
Me(Marple): I defined the existing API. It is not set in stone, so I'm open to changes. Also feel free to ask me about why things look the way they do.
Roman: Roman's process to resolve identities might also produce attribute (name, phone number, address, etc.). We should include him in review as a stakeholder (at least) since he'll be using the API to write.
Patrick: Patrick is working on the frontend, so we should include him in review as a stakeholder (at least) since he'll be using the API to read (and possibly for writes in some rare cases).