

```markdown
# SOVRA Sync – Meeting Notes

**Meeting title:** SOVRA Sync  
**Participants (selected):** Adam Bell, Alex Marple, Brendon Bluestein, David Dezso, Jimmy Qian, Lucy Love, Mitch De Jong, Sheil Mehta, Tim Jones, others  
**Context:** Weekly sync to align on OKRs, product development progress, data vendors, science/IP, Intel, marketing, and UX for SOVRA MVP and Texas pilot.

---

## 1. Admin & Housekeeping

- Some attendees having issues submitting expenses via **Justworks** (no expense tab).
  - Temporary solution: send expenses directly to **Adam**.
- Ensure **Doc** is added to the recurring meeting invite.
- Theme for **Friday** meetings:
  - Review “next week big buckets”
  - Confirm status vs. OKRs and product development
  - Mitch maintaining a **9-bucket tracker** by section and POC.

---

## 2. Product / Tech & Data

### 2.1 Identity Resolution & Tech

- **Identity resolution (Tim & Jimmy)**:
  - Large workstream; tracking toward an **MVP date ~ June 1* (loose tracking date).
- **Patent / IP work** (Tim, Jimmy, Doc):
  - Ongoing; multiple filings in flight.
  - Also tracked with a **June 1** target.

### 2.2 MVP PRD & Task Tracking (Tech)

- **Jimmy**:
  - Compiled offsite outcomes into an **MVP PRD**:
    - Clear **in-scope vs. out-of-scope** for July timeframe.
    - To be reviewed alongside **Sheil’s Figma** designs.
  - PRD to be shared in the **#all-sovra Slack channel**.
  - Request: team to review and provide feedback by **tomorrow** so scope can be locked.
- Task tracking:
  - Jimmy setting up a **more defined tech task tracker**.
  - Will be mapped into **Mitch’s 9-bucket tracker** for cross-team visibility.

### 2.3 Engineering & Mobile

- **Patrick**:
  - Exploring **mobile app development**:
    - Internal mobile-ready experience by **July**.
    - App stores target: **September**.
- **Nate**:
  - Supplementing team with **two additional engineers** (one mobile).
- **Alex Maeda**:
  - Supporting **data modeling and infrastructure**.
- **Mike (contractor)**:
  - Built a **data spreadsheet** for review by **Nate** and **Tim**.

### 2.4 Data Vendors & External Sources

- **PDL (People Data Labs)**:
  - Contract under review; questions from **Dan**.
  - Flat-file deal (~$280k upfront) postponed.
  - Requires Adam’s sign-off after questions clarified.
- **Pitbull**:
  - Contract signed, waiting on detailed **access** confirmation.
- **Trestle**:
  - Ready to go; just needs **credit card on file**.
- **Axle**:
  - Questions answered, awaiting response.
- **Anomaly Six**:
  - Secondary data source, likely **post-MVP**, out of current scope.
- **Carbon Arc**:
  - Appears strong on **commercial/consumption** data.
  - Less clear on **individual-level** data; under evaluation for gaps/replacements.

---

## 3. Go-To-Market: “Take Texas” & Market Research

### 3.1 Take Texas Go-To-Market Plan

- Lead: **Michael**, with **Laura** and **Gabby**; **Maddie** executing plan.
- Progress:
  - Good meeting with **Maddie**; moving toward detailed numbers for **Doc**.
  - **Target date:** **May 1** for a **75–90% plan** to feed Doc’s model.

### 3.2 Quantitative Survey

- Quant survey timing:
  - Vendor previously gave ~1 week, but ran longer with fewer responses last time.
  - Current expectation: **a couple of weeks**.
  - Provisional **target date: May 15**.
- Flexible reporting:
  - Ability to pull interim reports at **N=600, 700, etc.** to hit publishing dates without full completion.

### 3.3 Media & Influencers

- Need to:
  - Finalize **media costs** in **Austin, Houston, Dallas**.
  - Select **influencers** and estimate their costs.
- **Lucy** has a large influencer spreadsheet (including **Maggie’s** suggestions) to support planning.

---

## 4. UX Testing, Materials & Research

### 4.1 Discussion Guides & UX Testing

- **Sheil, Lucy, Michael**:
  - Writing **discussion guides** for UX testing.
  - Internal review this week; then team review.
- Dates:
  - Research begins by **Thursday**.
  - Guides and UX testing work will be in-field soon after.
  - Tracking **UX discussion guides completion by May 8**.
- **User testing materials package** (Sheil):
  - On track to finish by **Thursday**.
  - Tracking **May 1** as target for the full package.

### 4.2 UX Testing Strategy (Greta & Sheil)

- Approach:
  - Land users directly on dashboard and **observe natural behavior**.
  - Watch for **overwhelm**, navigation issues, and discoverability of key information.
- Tasks:
  - Always give users a **concrete task**, e.g.:
    - “Decide: would you date this person?”
    - “Find a specific data point” (e.g., county or state of a record).
  - Avoid leading questions; focus on:
    - Expectations vs. reality.
    - Whether they feel they successfully **completed the task**.
- Moderation:
  - Facilitate, but primarily **stay quiet and listen** to capture genuine reactions.

---

## 5. Dashboard & Figma: Self Search, Green Flag, Red Flag

### 5.1 Design Principles

- **Goal:** Dashboard that:
  - **Scales cleanly to mobile** (card-based, responsive).
  - Balances **data richness** with **approachability** for non-technical users.
- Collaboration:
  - Major work by **Sheil** and **Greta**, incorporating:
    - Offsite feedback.
    - **Doc’s/Tim’s** work on **three personality archetypes**:
      - **Green flag**
      - **Self search**
      - **Red flag**

### 5.2 Structure & Components

Key elements across archetypes:

- **Top nav + chat side panel**
  - Navigation sections with a vertical **scroll of cards**.
  - Chat can drive navigation (e.g., “show criminal records” auto-scroll).
- **SOVRA Assessment Overview**
  - High-level summary with key points “at a glance”.
- **Data leak / password card**
  - Users can **address/dismiss** specific leaks.
- **Marital status clarity**
  - E.g., explicit statement when no marriage records are found.
  - For red flag cases, clear **“likely married”** callouts with confidence score.
- **Social media footprint**
  - Reformatted for responsive display but conceptually similar to prior versions.
- **Public records**
  - Civil, criminal, marriage records, etc.
  - Chips/tags for **felony**, **violent**, or other seriousness indicators.
- **SOVRA “chime in” cards**
  - Injects interpretive or personality-oriented commentary.
- **In the news**
  - Articles with **sentiment** flags (positive/negative).
- **Geographic footprint**
  - Visualizing locations over time; still refining how many **sources** to visibly list.
- **Timeline**
  - Tailored narrative timelines per persona (self-search / green / red).
- **Final SOVRA summary**
  - Optional closing narrative at end of report.

### 5.3 Target Outcome

- Ready to begin **user testing by Thursday**.
- Main research questions:
  - Does the dashboard give **enough confidence** to decide whether to date someone?
  - What’s missing or hard to find?
  - Where are users confused or overwhelmed?

---

## 6. Science & IP

### 6.1 SBAM Prompt Engineering & Patterns

- Ongoing work with **Alex** on:
  - **SBAM prompt engineering**.
  - Drafting **checkbox outputs** for:
    - Green flag
    - Red flag
    - Personal-use “how to read” interpretations.
- **IP development**:
  - Reviewing prior **Bell’s use cases**.
  - Next major item: **Patent #2 – classification pipeline**:
    - To be developed with **Jeff** upon his return from honeymoon.

### 6.2 Raters & Calibration

- Five sites running **interaction ratings** and **use case ratings**.
- **Mitch + DWE flow**:
  - Calibration flow to shape rater behavior and outputs.
- Goal:
  - Move toward **two raters per subject** to cover more subjects overall.

### 6.3 fkT Patent

- **Geeta** has what she needs from the team.
- Awaiting additional input from **Jeff**.

---

## 7. Intel, Governance, Legal

- **sBIT patent submission**:
  - Updates being incorporated based on Mount Psych Houston discussions.
  - **Bell** reviewing and adding questions/recommendations.
  - Target: **package ready by end of week** to send back to **Jeff’s team**.
- Additional data companies:
  - Ongoing scan for other potential **data providers**.
  - Early impressions of **Carbon Arc**: stronger on **commercial transactions** than personal profiles.
- Governance / privacy:
  - **Minh** drafting initial **privacy policy** aligned with the **Texas pilot**.
  - Starting **data opt-out** conversations with tech:
    - Define at least **two contact channels** for opt-out.
    - Dedicated email for opt-outs already set up.

---

## 8. Marketing, Competitors, and UX Insights

- **Patrick** and **Lucy**:
  - Prior competitor reviews and intel work.
- **Laura**:
  - Deep dive on **dating-specific competitors**:
    - Exactly **how they present findings**.
    - **Gaps** in their UX and experience.
    - How they **develop and surface insights**.
  - Output will:
    - Inform **UX decisions**.
    - Feed into **go-to-market and media strategy**.
  - Next phase: analyze how these platforms drive **cultural conversation** and **go-to-market**.

- Deliverables:
  - **Michael** wants to:
    - Share findings with **David** first.
    - Then distribute to **the larger team** once complete.

---

## 9. Wrap-Up & Next Steps

- Mitch’s **9-bucket tracker** will be the core tool for:
  - Weekly status vs. OKRs.
  - Cross-functional visibility (tech, science, Intel, marketing, UX).
- **Sheil & Greta’s** updated dashboard designs:
  - Positively received.
  - Expected to ease **mobile** implementation and speed up front-end work.
- Onboarding:
  - **Sheil** to follow up with Adam regarding **Maddie** onboarding.
  - **Michael** to handle onboarding items for **Gabby** and **Maddie** separately.
- Meeting closes with:
  - Open offer from **Adam/Dave** for 1:1 follow-ups.
  - Plan to use upcoming Friday calls to refine **big buckets and priorities**.

---

## 10. Action Items

1. **Texas Go-To-Market & Quant**
   - Provide Doc with **quant results and numbers** for Texas GTM once available.  
     **Owner:** Michael  
   - Finalize **Take Texas GTM plan** to 75–90% by **May 1**.  
     **Owner:** Michael / Maddie / Laura / Gabby  
   - Run quant survey and aim for **May 15** as tentative completion date; send interim reports (e.g., at N=600, 700).  
     **Owner:** Lucy / vendor

2. **Science & IP**
   - Complete **SBAM prompt engineering** and checkbox drafts for green flag, red flag, and personal-use cases.  
     **Owner:** Science lead (Doc/Team)
   - Work with **Jeff** on **classification pipeline (Patent #2)** when he returns.  
     **Owner:** Science lead
   - Finish **sBIT patent** updates and send to **Jeff’s team** by end of week.  
     **Owner:** Intel lead / Bell
   - Support **fkT patent** as needed; Geeta to coordinate with Jeff.  
     **Owner:** Geeta / Science

3. **Data & Vendors**
   - Answer Dan’s questions and move **PDL** contract forward; Adam to sign if conditions met.  
     **Owner:** Jimmy / Adam  
   - Confirm and activate **Pitbull** data access.  
     **Owner:** Jimmy  
   - Activate **Trestle** (payment + usage).  
     **Owner:** Jimmy  
   - Evaluate **Carbon Arc** data for gap-filling or quality replacement; confirm whether it adds individual-level value.  
     **Owner:** Intel / Jimmy  
   - Continue exploring and prioritizing additional data companies (including Anomaly Six for post-MVP).  
     **Owner:** Intel / Nate / Jimmy

4. **Product & Tech**
   - Set **June 1** tracking dates for:
     - **Identity resolution MVP**.  
     - **Patent/IP work** (Tim/Jimmy/Doc).
     **Owner:** Mitch (tracker updates)
   - Add and track:
     - **May 1** for **Take Texas GTM**.  
     - **May 8** for **UX discussion guides / research**.  
     - **May 1** for **user testing materials package**.  
     **Owner:** Mitch
   - Share **MVP PRD** in #all-sovra Slack.  
     **Owner:** Jimmy  
   - Team to review PRD and provide feedback by **tomorrow**.  
     **Owner:** All relevant leads  
   - Set up a **more detailed tech task tracker** and align it with Mitch’s tracker.  
     **Owner:** Jimmy  
   - Progress mobile app plan for:
     - **July** internal-ready mobile experience.  
     - **September** app store release.  
     **Owner:** Patrick + mobile engineer

5. **UX & Research**
   - Finalize **UX discussion guides** and begin testing (target **May 8** completion).  
     **Owner:** Sheil / Lucy / Michael  
   - Finalize **user testing materials package** by **May 1**.  
     **Owner:** Sheil  
   - Conduct user tests:
     - Give users **clear tasks** (e.g., “decide whether to date”, “find a specific piece of info”).  
     - Observe navigation, overwhelm, and comprehension.  
     **Owner:** UX team (Greta / Sheil)

6. **Marketing & Competitors**
   - Complete and share **competitor UX/experience analysis**:
     - First to **Adam and David**.  
     - Then to the wider team.  
     **Owner:** Laura / Michael  

7. **Governance / Legal / Ops**
   - Continue drafting **privacy policy** for Texas pilot.  
     **Owner:** Minh  
   - Define and implement **data opt-out** flows (two contact channels, email live).  
     **Owner:** Minh + Tech  
   - Handle onboarding paperwork for **Gabby** and **Maddie**.  
     **Owner:** Michael  
   - Confirm any additional onboarding needs for **Maddie** with Sheil/Adam.  
     **Owner:** Sheil

---

```


