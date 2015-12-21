
Podcast manager: podbaby.me
===========================

User stories
------------

- As a site visitor, I can sign up and login to the app.

- As a (logged in) user, I can see all new episodes on the dashboard for podcasts I am subscribed to (nice to have: suggestions based on past choices) 

- As a user, I can search for new podcasts and subscribe to them. My search will also return podcasts I am subscribed to but these will be indicated.

- As a user, I can manually refresh my dashboard or podcast episode list (nice to have / do later: web socket that lets you know if new episodes eg. like twitter web app)

- As a user, I can unsubscribe from podcasts I am currently subscribed to. I will no longer see episodes for these podcasts in my dashboard.

- As a user, I can see all the episodes for a podcast and details about the podcast in its own detail view.

- As a user, I can stream an episode or download it. 

- As a user, I can manually mark an episode as being listened. It should still appear with some visual indicator e.g. faded out. It should no longer appear on my dashboard.

- As a user, I can mark a whole podcast (all episodes) as listened.

- As a user, I can see summary info on an episode.

- As a user, I can "star" episodes I'd like to watch later and see them in my dashboard under a "watch later" tab/filter.

- As a user, I can logout.

- As a user, I can deactivate/delete my account.

- As a user, I can have a reset link sent to my email address if I forget my password.

- As a user, I can can change my email address and password.

- As a user, I can see the most popular podcasts (as per total number of subscriptions).

- As a user, I can add a new podcast if the podcast does not already exist in the system by providing the RSS feed URL. I can then subscribe to the podcast. Once a podcast has been added others can find it in their search, but they should not see who added it.

- As a user, I can subscribe to a podcast. When I subscribe I can provide my own categories e.g. "comedy". If others have provided these they may be pre-selected.

Technology
----------

Exact technology TBD. Immediate choices for consideration:

Frontend
---------

- Babel/Webpack 
- React/Redux
- Bootstrap/CSSModules
- Mocha/Sinon/Chai

Backend
-------

- Go
- PostgreSQL
- Redis

DevOps
------

TBD, leaning towards Docker & Digital Ocean.

Data model
----------

ID PKs, timestamps implied.

User
----
Name (U)
Email (U)

Podcast
-------
URL (U)
Name
Description
Image
Website
User (FK -> User)

Subscription
------------
UserID (FK -> User)
PodcastID (FK -> Podcast)

Category
--------
Name (U)

SubscriptionCategory
--------------------
CategoryID (FK -> Category)
UserID (FK -> User)

Episode
-------
PodcastID (FK -> Podcast)
Title
Date
Description
MediaURL
MediaType

UserEpisode
---------
UserID (FK -> User)
EpisodeID (FK -> Episode)
IsListened
IsStarred

Screens
=======

Layout
    Logged out:
        About
        Contact us
        Login/Signup
    Logged in:
        Logout
        Settings
        Home (-> dashboard)
        Podcasts (-> list of all my podcasts)

Splash screen
About us (TBD)
Contact (TBD)
Login 
Signup
Recover password
Dashboard (show latest podcasts)
Search (show list of podcasts)
Podcast detail (shows list of episodes, description in panel)
Episode detail (shows single episode)
Settings (change email/password)


