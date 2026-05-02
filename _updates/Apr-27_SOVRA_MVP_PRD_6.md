**SOVRA MVP — Product Requirements Document**

**Status:** Draft v1.9 — for internal review

**Targets:** Closed beta — July 2026 · City launch (1–3 US cities — candidates Houston / Austin / Dallas) — September–October 2026

# **Key Terms**

A short reference for terms that recur throughout this document.

* **MVP —** Minimum Viable Product. The first version of SOVRA we ship to real users.

* **Subject —** the person being looked up in a SOVRA report (not the user doing the looking-up).

* **Self-search —** when a user looks up a report on themselves. Unlocks the Cybersecurity Leaks section, which is private to that user.

* **S-BAM —** Sovra Behavioral Analysis Model. The AI system that turns raw data about a subject into the trait ratings, risk assessments, and written narratives shown on the report.

* **Trait ratings —** S-BAM's 1–5 scoring of a subject across personality, personality-under-stress, and other behavioral dimensions, each with a confidence level and supporting evidence.

* **Tier —** a user's access level. Three buckets: Any (unverified visitor), Verified (free, ID-verified), Paid (verified \+ paying; multiple paid tiers).

* **Verified —** a user who has completed identity verification with our third-party identity provider. Required to view full subject reports and to access self-search Cybersecurity.

* **FCRA —** Fair Credit Reporting Act. US law governing background-check-style reports. SOVRA is NOT an FCRA-regulated product, and we display a disclaimer to that effect to prevent users from misusing reports for hiring or tenant screening.

* **HITL —** Human-in-the-Loop. People reviewing AI outputs to check quality and to label data that improves the model.

* **Beta —** the closed beta launch in July.

* **City launch —** the public launch in Sep/Oct to anyone living in 1–3 chosen US cities.

# **1\. Overview & MVP Goal**

SOVRA is a behavioral risk-report platform. A user searches for a person and receives a multi-section dashboard summarizing public-record signals, social-media presence, news coverage, geographic activity, life timeline, and (for self-search) cybersecurity exposure. A SOVRA Agent (chat) sits beside the dashboard and helps the user interpret findings.

The MVP serves a closed beta in July, then opens to one to three US cities in September-October. Beta users are early invitees with manual onboarding overhead acceptable; city launch is the first scaled-traffic event.

### **MVP success criteria**

* Beta (July): a user can complete the full happy path end-to-end (search → results → click into subject → review all sections → use Agent → upgrade to paid → run more searches).

* Beta: the SOVRA Assessment narratives, status tags (e.g. "Likely Married"), and confidence scores meet the quality bar set during HITL eval review.

* Beta: highly secure data infrastructure that prevents PII exposure and other data leakage.

* City launch (Sep/Oct): system handles steady-state traffic from one to three mid-size US metros. Search and report-generation latency stays within acceptable limits.

* City launch: paid conversion funnel is instrumented end-to-end and feeding analytics.

* City launch: native mobile app available with feature parity to web.

### **What is explicitly NOT in MVP**

* Subject continuous monitoring beyond ad-hoc "report ready" notifications (deferred entirely — see §6).

* B2B / enterprise features (concierge data services, bulk lookup, API access).

* Geographies outside the launch cities for the city-launch phase.

# **2\. Core User Flows**

### **2.1 Anonymous visitor**

* Lands on home page.

* Submits search query. Sees Search Results page with candidate cards (name, masked phone/email, location, age, gender, employer, education, breach count, alias tags, social-account count).

* Cannot click into a Subject Detail page. Clicking a result triggers a Create-Account / Login modal.

### **2.2 Verified (free) user**

* Account created with email \+ password; email and phone both verified via one-time codes.

* Identity verification via our third-party identity provider before first full report view (one-time).

* Can run searches, click into Subject Detail, view all dashboard sections at the Verified tier (see §3 matrix).

* Can self-search and view their own Cybersecurity Leaks section ("For your eyes only").

* Hits per-period report cap (TBD) → upsell to Paid.

### **2.3 Paid user**

* All Verified capabilities, plus tier-specific unlocks (more reports, full social handles, in-depth Agent — see §3 matrix).

* Subscribes through our payment provider's hosted checkout. Manages their subscription (upgrade, cancel, payment method) through the same provider's customer portal.

* Server-side enforcement of tier on every data-returning endpoint.

### **2.4 Returning user**

* "Welcome back, {first name}" landing.

* Resume previous queries from Search History (date-range \+ keyword filters).

* Receive Notifications when async report generation completes.

# **3\. Tier × Capability Matrix**

Tiers: Any (unauthenticated) → Verified (free, ID-verified) → Paid (Tier 1 / Tier 2 / Tier 3 — all require verification). Number of reports per period is TBD across Verified and Paid; lock before beta.

| Capability | Any | Verified | Paid T1 | Paid T2 | Paid T3 |
| :---- | :---- | :---- | :---- | :---- | :---- |
| Run search, view results page | ✓ (limited) | ✓ | ✓ | ✓ | ✓ |
| Open Subject Detail page | — | ✓ | ✓ | ✓ | ✓ |
| Overview / Relationship / Public Records / In the News / Geographic / Timelines | — | ✓ | ✓ | ✓ | ✓ |
| Social section — counts & summaries | — | ✓ | ✓ | ✓ | ✓ |
| Social section — full handles & deep links | — | — | ✓ | ✓ | ✓ |
| SOVRA Agent — short answers, common Qs | — | ✓ | ✓ | ✓ | ✓ |
| SOVRA Agent — in-depth analysis | — | — | ✓ | ✓ | ✓ |
| Reports per period | — | TBD | TBD | TBD | TBD (highest) |
| Search history retention | — | Indefinite | Indefinite | Indefinite | Indefinite |
| Cybersecurity Leaks — self-search only | — | ✓ (basic) | ✓ | ✓ | ✓ |
| Cybersecurity self — full PII detail (e.g. SSN, full breached PWs) | — | — | ✓ | ✓ | ✓ |

# **4\. Feature Specification**

## **4.1 Auth & Onboarding**

* **Account creation:** email \+ password, accept Terms of Service and Privacy Policy.

* **Email verification:** one-time code sent to the user's email.

* **Phone verification:** one-time code sent via SMS at signup.

* **ID verification:** third-party identity-verification flow that the user completes once. Required before viewing any full subject report.

* **Self-match:** a user is treated as "self-searching" when their verified identity matches the searched subject. This unlocks the Cybersecurity Leaks section, which is shown only to that user.

* **Sessions:** standard secure session management; users stay signed in across visits.

## **4.2 Search & Search Results**

### **Inputs**

* Free-text query ("Robert Schmidt, 56, NYC"). Optional structured filters at results level: Location, Age, Status, Current Company.

* Optional file upload (e.g. photo) for image-search (P2 — in MVP scope but lower priority than text search).

### **Output**

* Up to N candidate cards with: photo, name, masked phone, masked email, location, age, gender, occupation, employer, education, count of breach data records, count of social accounts, alias tags.

* Total result count visible in header ("Search results (4)").

### **Behavior**

* **Anonymous:** results page accessible. Clicking a candidate triggers Create-Account / Login modal (does NOT navigate).

* **Verified+:** clicking candidate navigates to Subject Detail.

* **Low-confidence multi-candidate matches:** UX exists for ambiguous resolution — show all plausible matches with confidence indicators.

* **FCRA disclaimer:** displayed before showing any results, persistent in-page.

## **4.3 Subject Detail Page**

Section anchors (left-to-right): Overview · Relationship · Public Records · Security · In the News · Social · Geographic · Timelines. Header includes back arrow.

Each section component supports four states: Populated · Empty (we don't have data) · Processing/Loading · Greyed-out (tier-gated). All sections render at the same time; loading states are independent.

The descriptions below specify what data each section is responsible for and what the user can do with it. Specific layouts, filter chip labels, category groupings, table column structures, and visual treatments shown in current Figma are directional only — final design will be fitted to the data we actually have available, in collaboration with design.

### **4.3.1 Hero**

* Identifying information for the subject: name, age, primary location, photo.

* At-a-glance summary signals from across the report (e.g. counts of contributing sources, presence of significant findings such as cyber exposure).

* Topline status indicators for the highest-level questions a user is likely to ask first (e.g. relationship status, criminal-record presence).

* "Key Areas" pointers — deep links to the 1–3 sections most relevant for this subject.

### **4.3.2 Executive Summary (behavioral overview)**

* Sits near the top of the report, before the data sections. Rendered as a prominent narrative block in the SOVRA-Assessment style.

* Content: 3–5 sentence S-BAM-generated behavioral overview of the subject. Frames who the person appears to be, top behavioral traits, dominant risk profile, and the 1–3 sections that warrant the user's closest attention (cross-references the Hero "Key Areas" pointers).

* Always present in MVP (P0). Default behavior if S-BAM cannot produce an overview at acceptable confidence: hide the block rather than render a low-confidence summary.

### **4.3.3 Relationship Status**

* Determination of relationship state with associated confidence (specific status labels and confidence representation are design choices).

* Supporting counts and signal types that contributed to the determination.

* Narrative paragraph explaining the reasoning.

* Specific evidence shows data counts and source types, not record-level detail.

### **4.3.4 Public Records**

* Records pulled from public sources, organized into categories and surfaced for the user to scan.

* Specific category breakdown, filter affordances, and ordering are design choices — final categorization will be fitted to what coverage we have in the launch cities.

* Each populated record shows: type, brief description, date, jurisdiction. Empty categories are surfaced with explicit "none found" affordances so the user can distinguish absence from missing data.

* Categories in MVP coverage scope: Civil, Criminal, Marriage, Sex Offender Registry, OFAC. Other record categories deferred to post-MVP.

### **4.3.5 Security / Cybersecurity Leaks**

* Self-search only (verified user, subject \= self). Section visually marked as private to the viewer.

* Surfaces credential and account exposure detected in breach data and on the open / dark web.

* Tier-gated detail: Verified users see exposure presence and high-level account references; Paid users see fuller PII detail (e.g. full breached passwords, SSN-class data) under explicit reveal interactions.

* Includes guidance on remediation steps (P1).

* SOVRA Assessment narrative summarizing exposure level and recommended actions.

### **4.3.6 In the News**

* Articles and news mentions associated with the subject.

* Filtering capabilities available to narrow the list.

* Per-article display: article metadata, sentiment indicator (if classifier reliable), and a way to read or open the source.

* Sentiment classification per article is P3 — ship sentiment indicator in beta if model confidence is acceptable; otherwise hide and ship by city launch.

### **4.3.7 Social**

* Subject's presence across the social platforms in MVP scope: LinkedIn, Instagram, Facebook, TikTok, Twitter.

* Per-platform: presence indication, account counts, follower / following / connection metrics, posts or articles visible.

* Per-account drill-down with confidence indicator and a link out to the source.

* Sample of public content (posts, photos) where available, with timestamps.

* Tier gating: account counts and summaries available to Verified; full handles, deep links, and richer per-account detail are paid-only.

### **4.3.8 Geographic**

* Where the subject has been associated with locations, derived from social posts, news mentions, and records.

* Visual representation of geographic activity on a map. Beta ships with binary location markers; weighted heatmap follows by city launch.

* Drill-down: per-location, the source records that tie the subject to that place.

* Recency emphasis — most recent / current location surfaced prominently.

### **4.3.9 Timelines**

* Chronological view of life events derived from the subject's records.

* Event categories in MVP: Professional, Education, Personal, Location. Specific lane organization, color encoding, and time-axis affordances are design choices.

* Each event is interactive — clicking surfaces the supporting evidence behind it (P1, required for MVP).

* Date-range scoping for navigating long timelines.

* S-BAM patterns / Inflection Points overlay: P3, post-MVP.

### **4.3.10 SOVRA Assessment blocks (behavioral analysis)**

Each major section of the report renders a SOVRA Assessment block. Generated by S-BAM, each block has two responsibilities working together:

* Findings summary — concise 1–2 sentence recap of what the section's data shows.

* Behavioral analysis — 1–3 sentences interpreting what those findings indicate about the subject's behavioral profile, tying to S-BAM's relevant trait ratings or risk categories. For example, a Public Records block might note that a pattern of civil disputes with no criminal record points to elevated reputational risk but low safety risk; a Social block might note that an active, public-facing presence indicates extroversion with no flags for harassment patterns.

* **Required in MVP for:** Relationship Status, Public Records, Security (self-search), In the News, Social. Geographic and Timelines do not require blocks for MVP — ship if S-BAM coverage permits, otherwise omit.

* **Visual treatment:** consistent, recognizable block styling across all sections. Specific styling is a design decision.

* **Quality gate:** blocks must read as coherent and defensible in HITL review on the test cohort before they ship. Bad behavioral analysis is worse than no behavioral analysis — a misfire here directly damages user trust in the product.

* **Source:** S-BAM section narratives, conditioned on the S-BAM trait and risk outputs (see §5.2).

## **4.4 SOVRA Agent (chat)**

* Chat panel renders to the left of the Search Results and Subject Detail pages.

* Initial agent message after a search: confirms the search, updates the search result panel.

* Agent can answer scoped questions about the current subject.

* **Guardrails (P1, hard requirement):** Agent must NOT give hiring advice, credit advice, clinical psychiatric diagnoses, or return raw PII. Must reject prompts that attempt to bypass system prompts.

* Async report generation: long-running analyses produce a notification when ready (do not block chat).

## **4.5 Notifications**

* Three categories: All · Security · Reports.

* Each notification: icon, title ("Jose Campos report ready"), subtitle ("Dashboard is now ready to be analyzed"), relative time, action buttons (Open this report / Dismiss).

* Day grouping: Today / Yesterday / Older.

* Mark-all-as-read action.

* Notifications load progressively (oldest entries fetched as the user scrolls).

* Push notifications for report-ready alerts — deferred to city launch unless straightforward to add earlier.

## **4.6 Search History**

* Keyword search across past queries.

* Date-range filter (with quick "Edit" affordance).

* Day grouping: Today / Yesterday / weekday-dated entries.

* Each row: timestamp · query text · count of records returned for that query.

* Server-side pagination, infinite scroll.

* Search history is retained indefinitely for all users with an account.

## **4.7 Settings**

* Preferences tab — light/dark theme, notification preferences. Other onboarding-driven preferences gate UI elements.

* Account tab — email, password reset.

* Subscription tab — current tier, link to manage subscription (handled by our payment provider).

* Verification tab — current ID-verification status, option to re-run.

## **4.8 Pricing & Tier-Gating**

* Pricing page: three Paid tiers (Tier 1 / Tier 2 / Tier 3 — final naming TBD) plus a free Verified tier. Monthly billing for MVP; annual billing is deferred.

* Payment processing handled by an industry-standard provider; users manage their subscription (upgrade, downgrade, cancel, payment method) through the provider's hosted portal — we don't store credit card details.

* Tier checks happen on the server for every piece of data we return, so the user's tier can't be bypassed by manipulating the front-end. UI greying-out of higher-tier features is informational only — the security boundary lives on the server.

## **4.9 Concierge**

* At launch, Concierge is a contact point for users — a way to reach a human at SOVRA for help, questions, or escalations.

* Stated capabilities (what users are told they can request from Concierge) are TBD.

# **5\. What S-BAM Provides to the Report**

S-BAM (the Sovra Behavioral Analysis Model) is the AI component responsible for behavioral ratings and free-form behavioral narratives in the report. Other report content — relationship status, criminal record presence, public records, life events, geographic activity, social presence counts, and the like — is derived programmatically from source data and is out of S-BAM's scope.

### **5.1 What S-BAM consumes**

* All available data we've gathered about a subject: public social-media activity, court and other public records, news coverage, web-search results from approved sources, and breach data.

* Every piece of data carries a record of where and when it came from, so any output S-BAM produces can be traced back to its source.

### **5.2 What S-BAM produces for the report**

* Trait ratings — 1–5 scoring of the subject across personality, personality-under-stress, and other behavioral dimensions, each with a confidence level and supporting evidence.

* Executive Summary narrative — 3–5 sentences summarizing the subject's behavioral profile and which sections of the report deserve closest attention. Powers the report's Executive Summary block (§4.3.2).

* Section narratives with behavioral analysis — for each section that gets a SOVRA Assessment block (§4.3.10): a short findings summary plus a behavioral interpretation grounded in S-BAM's trait outputs. Required in MVP for: Relationship, Public Records, Security, In the News, Social.

### **5.3 Guardrails**

* S-BAM never produces hiring or credit advice, never gives clinical diagnoses, and never returns sensitive personal information that hasn't been explicitly approved for the relevant tier.

* S-BAM is hardened against attempts by users to manipulate its instructions through crafted input.

* S-BAM quality is measured continuously through human review and automated quality scoring; no model change reaches users without passing regression checks.

# **6\. Out of Scope for MVP / Explicit Deferrals**

* Subject Monitoring (saved subjects with delta alerts) — fully out of MVP. Revisit post city-launch.

* Sentiment-based news filter — ship if classifier is reliable, otherwise hide chip until city launch.

* Geographic heatmap weighting — ship binary markers in beta, weighted heatmap in city launch.

* Mutual contacts within Social — high cost, low value; deferred.

* S-PACE / Pattern-of-Life / Inflection Points overlays in Timelines — P3.

* B2B / API access; expansion beyond the 1–3 launch cities.

# **7\. Open Questions**

* Number of reports per period — for Verified, Tier 1, Tier 2, Tier 3\.

* Final naming for paid tiers (currently Tier 1 / 2 / 3 placeholders).

* Final pricing for paid tiers.

* Onboarding-driven preferences — exhaustive list, and which gate which UI elements.

* Final decision on the launch-city count (1, 2, or 3).

* Concierge stated capabilities — what users are told they can request from a SOVRA human contact at launch.