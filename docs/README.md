# [@MyPackBot](https://t.me/MyPackBot) [![discord](https://discordapp.com/api/guilds/208605007744860163/widget.png)](https://discord.gg/KYQB9FR)

[![License](https://img.shields.io/crates/l/rustc-serialize.svg)](LICENSE)
[![Build Status](https://travis-ci.org/toby3d/MyPackBot.svg)](https://travis-ci.org/toby3d/MyPackBot)
[![Go Report](https://goreportcard.com/badge/github.com/toby3d/MyPackBot)](https://goreportcard.com/report/github.com/toby3d/MyPackBot)
[![Release](https://img.shields.io/github/release/toby3d/MyPackBot.svg)](https://github.com/toby3d/MyPackBot/releases/latest)
[![Patreon](https://img.shields.io/badge/support-patreon-E6461A.svg?maxAge=2592000)](https://www.patreon.com/toby3d)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftoby3d%2FMyPackBot.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftoby3d%2FMyPackBot?ref=badge_shield)

![bot logo](https://raw.githubusercontent.com/toby3d/MyPackBot/gh-pages/static/social/og-image.jpg)

## Wat?
This is a Telegram-bot that collects all the stickers sent to it in one (almost) infinite pack. No more, no less.

**Benefits:**
- Does not require creation of a set with a unique URL and/or name;
- Indeed (almost) unlimited pack size;
- Keeps stickers belonging to their original sets;
- Fully support the standard functionality of Telegram stickers (for example "add to favorites");
- Avaliable anywhere in Telegram by typing `@MyPackBot ` in the input field;
- Supports filtering of results by emoji's: `@MyPackBot 😀👍`;
- Fast as f\*\*\*king Sonic;
- Worked with uploadable WebP stickers;
- Worked with blocked by rightholders sets (but this is not exact);

**Disadvantages:**
- Requires type `@MyPackBot ` in the input field;
- Availability depends on the internet connection and bot uptime;
- Supports search/filtering only for first emoji associated with sticker;
- Does not support synchronization of the updated original set contents with the saved set contents in the bot.

## Why?
Because Telegram native tools for managing stickers are somewhat limited:
- User can have only **up to 200 active stickers sets**;
- In one set can be uploaded **up to 120 stickers**;
- User can have only **up to 5 favorites stickers**;

Having done simple mathematical calculations, we can assume that the **maximum user capacity** (when he has the maximum number of sets, each of which contains the maximum number of stickers) **is equal 24,000 stickers**.

But, as usual, there are problems:
- **Most of the sets are incomplete** and contain less than 120 stickers (sometimes - only 1-3 stickers on whole set);
- **Some sets contains junk, duplicated and promotional stickers**;
- **Sometimes user want use own stickers** by uploading WebP files, but without creating new sticker set;
- Anyway, **user just want have as many stickers as he want**;

To solve these problems, this bot was designed.

## How?
### tl;dr
- Telegram API [supports stickers as results](https://core.telegram.org/bots/api#inlinequeryresultcachedsticker) in inline query;
- Telegram API allows to use someone else's FileID for results;
- It is not necessary to [create a new set](https://core.telegram.org/bots/api#createnewstickerset) using Telegram, since it only "references" existing files;
- Bot saves only [user](https://core.telegram.org/bots/api#user) info, [sticker and name of his set](https://core.telegram.org/bots/api#sticker) in the database if user upload custom sticker or send/forward already existing;
- Database architecture allows to filter keys by user ID and sort them by set name and emoji value;
- When requesting inline query, bot simply create results from filtered database keys;
- ???????
- PROFIT!!1

## Step-by-step
I'm too lazy to write, so just check the source code for the comments. 👀

### Dependencies
Bot uses the following dependencies:
- Written on [Go](https://github.com/golang/go) language, because I <3 Go;
- I ventured to migrate to my own [telegram](https://github.com/toby3d/telegram) package to win in convenience and productivity;
- I use [dlog](https://github.com/kirillDanshin/dlog) for debugging without spamming on production server by use only one build flag;
- Data of users and stickers save thanks to [BuntDB](https://github.com/tidwall/buntdb);

## Support
### GitHub
You can [request fix/add some things](https://github.com/toby3d/MyPackBot/issues/new), [make a patch](https://github.com/toby3d/MyPackBot/compare) or help with [translation and localization](https://github.com/toby3d/MyPackBot/tree/develop/translations) on your language.

Ah, and star this repo, of course.

### Patreon
**I work on my own projects in my free time.** Please think about the [financial support for my independence](https://patreon.com/toby3d) so that I can devote more time to this bot and other projects. In exchange for an award!

### Social
Subscribe, follow my resources and feel free to maintain contact with me: https://toby3d.github.io

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftoby3d%2FMyPackBot.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftoby3d%2FMyPackBot?ref=badge_large)