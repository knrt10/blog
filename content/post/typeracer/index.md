---
title: ShellRacer
publishDate: 2018-05-14
draft: false
summary: Why I created a typeracing application.
draft: false
authors:
- admin
tags:
- js
- opensource
- node.js
- touchtyping
categories:
- opensource
image:
  placement: 2
  caption: "[Image](https://user-images.githubusercontent.com/24803604/39766933-2386f1d8-5303-11e8-9f34-76c3b53f58c7.png)"
  focal_point: "Center"
  preview_only: false
  alt_text: Typing practice from your terminal and features like practice mode and online competition mode.
---

## A little background

For months my eyes were set on one goal : the prestigious and famous Google Summer of Code. I had been contributing to my organization like crazy, creating PRs after PRs, raising issues and helping out fellow contributors and other aspiring GSoCers on the community's slack team.

The day of the results came and out of the two slots given to my organization my project wasn't in the list of selected projects. It was hard at first, I thought I had done everything right but still my name wasn't up there.

The next morning I decided not to waste my time anymore and started to look at the bright side of things, contributing to open source project helped me learn invaluable skills like _test driven development, CI/CD, git (rebasing, merge conflicts), etc_. I may not be getting any stipend for writing code throughout the summer but I earned some invaluable skills which I can use in my future projects.

So it was exam time and I was hacking on one of my project’s as usual (a CLI tool to upload images to a cloud service) when one of my cool friends ([Palash Nigam](https://github.com/palash25)) entered my room and said with a tone of surprise _“You don’t know how to touch type? That’s pathetic”_ :grinning:. So he introduced me to this site called typeracer (so I listened to him and started practicing on that site and I was hooked). One of my seniors from college ([Shibasis Patel](https://github.com/shibasisp)) gave me this idea of creating a CLI tool to play typeracer (as I had earlier asked him for some project ideas) which we could use to introduce the freshmen to both the shell and touch typing. This gave me a new purpose, so I started coding. After about a week of writing code I am proud to present to you **[typeracer-CLI](https://github.com/p-society/typeracer-cli)**

## What is typeracer-cli?

So it is basically a terminal client for playing typeracer on your shell. As soon as you start the game you will be presented with a paragraph which you have to type out and at the end your time and speed (in wpm) are recorded and presented as an output.

> What’s new about this? CLI versions of this game already exist.
Agreed, other versions of this already exists but they don’t offer all the features that this client does, like:

- Practice mode (offline mode)
- User stats (words per minute, time taken)
- Online mode (have a type-race by spawning up a server and sharing it with your friends)
- Ask for a rematch after the race ends (online mode)
- View the top 10 High scores in online mode

## The motivation behind it

Well, they say failure is a great motivator. I learned this the hard way. Failing in getting selected for GSoC motivated me even more to be a better developer. The other big reason was the desire to do something for my college. I had realized the benefits of working in a developer community during the time I was preparing for GSoC. Although I was a part of my college’s Programming Society I hadn’t contributed as actively as I should have. This project turned out to be one of the ways of contributing to my community by spreading awareness about touch typing and CLI tools among young aspiring developers who are just starting out.

## Implementation

Initially the task was to get keystrokes from the user’s terminal which at that time I thought was impossible. But I found about readline and keyspress events in nodejs which helped me to move further in coding.

The tasks were broken up into the following:

- Conver this tool to an npm package
- Offline practice mode
- Generate random paragraph for every race
- Add more sensible paragraphs
- Display the user’s time and speed as they type
- Setup server for online mode
- Improve the API Design
- Write tests

## Getting into every point in detail

### Converting it to an npm package

This was important task so that one can easily download the package and install it globally from [npm](http://npmjs.com/package/typeracer-cli). For that we need to use a very important line on the start of the file that is going to execute.

> #!/usr/bin/env node is an instance of a [shebang](https://en.wikipedia.org/wiki/Shebang_(Unix)) line: the very first line in an executable plain-text file on Unix-like platforms that tells the system what interpreter to pass that file to for execution, via the command line following the magic #! prefix (called shebang)

Although **Windows does not support shebang lines**, so they’re effectively ignored there; on Windows it is solely a given file’s filename extension that determines what executable will interpret it. **However, you still need them in the context of npm**.

### Implementing the offline (practice) mode

Initially some commands were written for the execution of practice mode. With the help of a package [commander](https://www.npmjs.com/package/commander) I was able to achieve this task.


```js
program 
 .command('practice')
 .alias('p')  
 .description('Starts typeracer in practice mode') 
 .action(() => 
   {  
     game() 
   })  
```

### game() function

This is main logic that allows the application to get keystrokes from the client, but we also have to listen to keypress event for completion of this task. `stdin.on('keypress', keypress)`

```js
const stdin = process.stdin
const stdout = process.stdout
stdin.setRawMode(true)
stdin.resume()
require('readline').emitKeypressEvents(stdin)
```

Now in game() I was enabling `keypress` event after 5 seconds of game, and showing paragraphs to user so that they get time to _relax their fingers, twist turn their neck, crack their knuckles and say “bring it on”_.

I was displaying three things to client when they were typing

- Real time analysis of their typing with green, red representing correct and wrong characters respectively.

```js
/**
* @function color
* @param {String} quote
* @param {String} stringTyped
*/

function color (quote, stringTyped) {
  let colouredString = ''
  let wrongInput = false

  const quoteLetters = quote.split('')
  const typedLetters = stringTyped.split('')
  for (let i = 0; i < typedLetters.length; i++) {
    // if a single mistake,
    // the rest of the coloured string will appear red
    if (wrongInput) {
      colouredString += chalk.bgRed(quoteLetters[i])
      continue
    }

    if (typedLetters[i] === quoteLetters[i]) {
      wrongInput = false
      colouredString += chalk.green(quoteLetters[i])
      if (quote === stringTyped) {
        gameEnd = true
      }
    } else {
      wrongInput = true
      colouredString += chalk.bgRed(quoteLetters[i])
    }
  }
  return colouredString
}
```

- Real time analysis of their speed in words per minute.

The following snippet explains how to get the speed of a user according to correct words typed by them.

```js
/**
* @function updateWpm
*/

function updateWpm () {
  if (stringTyped.length > 0) {
    wordsPermin = stringTyped.split(' ').length / (time / 60)
  }
}
```

- Calculating the time taken

```js
/**
* @function Time
*/

function Time () {
  time = (Date.now() - timeStarted) / 1000
}
```

In the end there is an option to retry where you can restart the match with generation of new paragraph every time.

## Online Mode

This was very important feature to implement as this sets this client apart from other CLI versions.
This was implemented using **socket.io** accordingly:

- Creating a server
- Connecting clients to the server
- Create private room for competition
- Send scores to all clients at end of game
- Rematch feature
- Random paragraphs for every race
- Top 10 high scores

## Getting into every point in detail

### Creating a server

I am quite fluent with Javascript and NodeJs, so I used NodeJs for creating server of the application and hosted in on [Glitch](http://glitch.com/).

[Socket.io](http://socket.io/) was used to provide web sockets for clients to connect and emit events for server to listen. [MongoDB](https://www.mongodb.com/) was used as database for storing top 10 high scores of clients.

### Connecting clients to sever

Initially the client part was quite tricky as I had not understood socket.io upto basic level. At first I was working on local server or you can localhost. I was using [socket.io client](https://github.com/socketio/socket.io-client) for client side but still took me a whole day to understand the basic and connect a client to the server.

### Creating a private room for competition

Now in socket.io you can create different namespace or you can say rooms to join. So I had to get some input from user to create a private channel where they can race otherwise it would create countless problems. I used the crypto node module to provide cryptographic functionality that includes a set of wrappers for OpenSSL’s hash, HMAC, cipher, decipher, sign, and verify functions and generated random strings with it every time.

```js
const roomNumber = crypto
            .randomBytes(12)
            .toString('base64')
            .replace(/[+/=]+/g, '')
```

Another problem was how the server would know the number of players in a room so that it can emit an event for race to start. For that I asked number of players from client they wanted to race with (using the npm package [inquirer](https://www.npmjs.com/package/inquirer).)

So when the user joined the room I emitted all the information of client to sever so that it can work according to that.

```js
// Emitting client info on joining the room

  _socket.on('room', function (val) {
    _socket.emit('join', {roomName: val.value, username: data.username, number: data.number, randomNumber: data.randomNumber})
  })
```

Now when the server knew that all the clients have joined the race it emitted an event to clients and started the race. It was important to send random paragraphs on every connection and also same paragraph to all clients in a same room.

```js
/**
* @function randomNumRetry
*/

function randomNumRetry () {
  randomNumber = Math.floor((Math.random() * paras.length))
  quote = paras[randomNumber].para
  if (quote.length < 100) {
    quote = paras[randomNumber].para + ' ' + paras[randomNumber - 1].para
  }
  return quote
}
```

When everything was in order and every client completed the race, the server emitted an event sending all the scores to every client and asking if they wanted a rematch. Similarly for a rematch random paragraph was generated.

### Top 10 high scores

For this feature to work I initially created a shell database with ten anonymous users with their scores initialized to 0. Now whenever someone plays an online game and score greater than 10th highest score in database, it replaces the 10th highest scorer with the user in the database (this was done to avoid excessive use of database).

```js
// Getting documents from databse

      Score.findOne({_id: process.env.ID}, (err, players) => {
        if (err) throw new Error(err)
        let playersArray = players.players.sort(function (a, b) {
          return b.score - a.score
        })
        let lowestScore = []
        lowestScore.push(playersArray[playersArray.length - 1].score)

        // checking if last score is less then current score
        function remove () {
          // First removing last player
          Score.update({_id: process.env.ID}, {$pop: {players: 1}}, (err) => {
            if (err) throw new Error(err)
            console.log('Removed last player')
          })
        }

        function add () {
          // Then updating current player
          Score.update({_id: process.env.ID}, {$push: {players: {score, username}}}, (err) => {
            if (err) throw new Error(err)
            console.log('Added new High score')
          })
        }

        async function update () {
          // Then again sorting it correctly
          await Score.update({_id: process.env.ID}, {$push: {players: {$each: [], $sort: -1}}}, (err) => {
            if (err) throw new Error(err)
            console.log('Sorted in descending order after adding')
          })
        }
        if (score > lowestScore[0]) {
          (async () => {
            Promise.all([update()]).then(async () => {
              await remove()
              await add()
              await update()
            })
          })()
        }
      })
```

## Support Us

This project was a great learning experience for me and we (my friends and I) are looking to build more such awesome projects in the future. We are a bunch of undergrads passionate about software development looking to make cool stuff. A little motivation and support helps us a lot. If you like this nifty hack you can support us by doing any (or all 😉 ) of the following:

- ⭐️ Star us on [Github](https://github.com/p-society/typeracer-cli) and make it trend so that other people can know about our project.
- Install it and increase our download count on npm.
- Tweet about it (our handle is psociiit).

Thanks to Palash Nigam for helping me to write this article and also Shibasis Patel for sharing this cool idea.

### Did you find this page helpful? Consider sharing it 🙌
