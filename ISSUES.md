 Issues reported in the latest (2025 dec. test).
 
 Im to lazy, so i made an FIXED list, instead of editing in the issue desc.
 Maybe later I will go over it line by line, and remove the fixed ones.

FIXED:
Cooking. When the player was used fancy, only fancy got removed from the inventory.
Less chance to win in muffin minigame. 1/12
Fixed the fortune wheel
 
 IS FIXED? = NO
 LEVEL = EASY
 - FORTUNE WHELL
   - ISSUES:
     - The fortune wheel is OP, players can farm 800+ golds in 2 days, and can spin the wheel 3000 times in 2 days.
     - Sometimes the decors not appaears in the inventory.
     - Players can win marketplace items, or different event items.
   - FIXES:
     - Make it to less chance to win gold.
     - Fix the inv.
     - Add a check to the items if theres a cash / gold value to the object (if not, then its not buyable, so marketplace object), and check the object event.

IS FIXED? = NO
LEVEL = EASY
- WAITERS SERVING
  - ISSUES:
    - When the counter and the chair is 1 empty space away, the waiters serve invisible foods. Its like not picking up the food from the counter.
    - Sometimes multiple waiters can serve one customer (or maybe with clean too)
- FIXES:
    - Make it to force wait 1s? (Or IDK), this issue is comming from deep!
    - Fix it.


IS FIXED? = NO 
LEVEL = BABY
- MASTERY LEVELS
  - ISSUES:
  - Player said: Mastery level progress seems to revert to an older number sometimes after reconnecting. So the mastery values not updates into the DB.
  - And force check if were using the mastery values.
  
- FIXES:
  - Save to the db
  - Modiy it to use mastery values.
  
IS FIXED? = NO
LEVEL = EASY
- ACHIEVEMNTS
  - ISSUES:
  - Some achievements not updating (not yet added)
  - On login, the server crashes, because when offline players login back again, and were adding cash to achievement, and when the level reached, its nill, so we cant send achievement earn.
  
- FIXES:
  - Implement all achievements.
  - Make it to bind the AchievementEarn to the player before LBU, and after we can send it.
  
IS FIXED? = NO
LEVEL = MEDIUM
- BUYING WALLS:
  - ISSUES:
  - Sometimes the wall is not updating, and sometimes its giving back incorrect ammount to the inventory. 
  - USER MEDIA: https://cdn.discordapp.com/attachments/1453807578936115414/1454129330237083834/Screen_Recording_2025-12-26_100453.mp4?ex=69573738&is=6955e5b8&hm=4330f693eb9544978a11b49fc44b6cd9b772a5fd8f72b4b27bd7f44faa551617&
  
- FIXES:
  - Rewrote the buying walls (EASIER THAN FIXING IT!)
  

IS FIXED? = NO
LEVEL = BABY
- REGISTER
  - ISSUES:
  - Allow more special characters, because players cant use 1-2 special characters as email, and cant use . in the email user.
  
- FIXES:
  - Allow more special characters
  - Allow . special character


IS FIXED? = NO
LEVEL = HARD
- CAFE EXPAND
  - ISSUES: Not yet implemented, but sometimes its can break the cafe.

- FIXES:
  - DO IT.

IS FIXED? = NO
LEVEL = EASY
- INSTANT COOKING
  - ISSUE:
  - Were using the dish duration (from the items) to check how much gold does it cost.
  - The players cant instant cook, but they have X Remaing instant cookings
  
- FIXES
  - Get remaing time in hours, and use gold from the remaing time, not the actual dish duration
  - Dont save it to the DB, calculate the max instantcooks on login and from the level ratio.

- ME:
heres the issue:
    if c.Player.GetGold() < dishCostGoldPerHours || c.Player.GetInstantCookings() > c.Player.GetMaxInstantCookings()


youre having enough money, but the maxInstantCookings value is lower than the instantCookings. 

MAX: 12
Youre having: 16

IS FIXED = NO
LEVEL = EASY
- COOP LEAVE
  - ISSUE:
  - Players cant leave coop.

FIXES:
  - Check it why...

IS FIXED = NO
LEVEL = BABY
- CAFE LUXURY
  - ISSUE:
  - Out of sync
  
FIXES:
  - Dont store the luxury in the db. Calculate it from the cafe objects, and set it.

  
CURRIER

when u visit someone, bots are eating all the time. if they take the whole room- new customers arent happy about it (offline player)

won a snowman decoration from the prize wheel, but didn't receive it in inventory

way to high chance for 10 gold in wheel 

won this guy and i dont have him in inventory
https://cdn.discordapp.com/attachments/1453807578936115414/1454762080136859648/image.png?ex=69578a44&is=695638c4&hm=83a6f3b432a365d3ca53106aa04d65c622a1c45c82f389daf07ef031d67491a3&

cant sell fridge, then when i leave from build mode she walks with my cursor in cafe
https://cdn.discordapp.com/attachments/1453807578936115414/1454765355653402781/image.png?ex=69578d51&is=69563bd1&hm=57601477ab9ed8ce56adceff03abc308de8a8646cce81f7480a5e4707cda8e47&
https://cdn.discordapp.com/attachments/1453807578936115414/1454765509152215081/image.png?ex=69578d75&is=69563bf5&hm=e38b702135b86880717531f391b7c39c989e19a1ff6df68342c704a7d75bc5f2&

Players won empty decoration
https://cdn.discordapp.com/attachments/1453807578936115414/1454765859552890942/image.png?ex=69578dc9&is=69563c49&hm=569500e6fd31fb7454198a1385bd093e6f6e964007a3a41b491c9ac792cf3b5d&

presents ratio on wheel - 3 or 2 when it always was 100/200 aswell visual bug - its still on my counter
https://cdn.discordapp.com/attachments/1453807578936115414/1454768344531210384/image.png?ex=69579019&is=69563e99&hm=05034476c9b7c43fc59f6a88ced744ae96d5cfcfe0af4ff4cd84e62167a5b547&

problem after receiving food from the gift after winning from the wheel
https://cdn.discordapp.com/attachments/1453807578936115414/1454823946468986880/image.png?ex=6957c3e2&is=69567262&hm=6d378c3c015a917cedb17cf4dea4aafdf8a6c8c69a06fc55ddbe05f57e258631&

when i click to kick player from cafe he is still inside
players can join other cafes while the owner in the editor mode

I bought a wall color for 5k but right in that moment I was disconnected. The money is gone, the wall too, in my inventory I got that default blue wall that’s still in the café.
Then I tried buying a second wall. When I decided to put the default blue above it so I can put the wall somewhere else I still got the blue wall (by now two of them) in my inventory but not the previously bought wall.


got 8 or something gifted and after a crash its at 2500
https://cdn.discordapp.com/attachments/1453807578936115414/1454952412153516207/image.png?ex=695792c6&is=69564146&hm=e88d779e850be1f32be9d98133c534f15b7466aa9c2e79eb1c1ccf847e0b94ce&
https://cdn.discordapp.com/attachments/1453807578936115414/1454952811723882659/image.png?ex=69579326&is=695641a6&hm=940d8535d66b7a12fdeabf1f4f421a5388e3808fa8669f7b3f3f119334b74ad0&

I found another problem with gifted food. In my case it visually showed food on the counter after it disconnected but it said empty counter. When I clicked on it it asked me if I wanna throw it out…

WRONG
Got the M'n'Cs in gift box from wheel.
Also when I put out these they were duped twice. However did this didn't happen to other dishes I put out from gift box. They were actually used. Maybe it takes a while for game to save. I guess staying the game longer or being offline longer makes sure to get rid of them. Also you can't throw anything. 
Oh wait NOW they disappeared they like they should. Game got synced at the end. Just takes a whiles
https://cdn.discordapp.com/attachments/1453807578936115414/1455165293797048452/image.png?ex=6957b049&is=69565ec9&hm=312ba89ba4f16e75022ccfe02235b096ef4e3895b64f26b148b0aae68f9c0fe5&

https://cdn.discordapp.com/attachments/1453807578936115414/1455167513947017256/image.png?ex=6957b25b&is=695660db&hm=435f23407d560d2205f8f08afd991e275a159d56b2c3ab5210d3284403d86ddd&
https://cdn.discordapp.com/attachments/1453807578936115414/1455167571275022346/image.png?ex=6957b268&is=695660e8&hm=15f6a01b6fc61eee11a0b4e308e64eca7992a169ab77ab7fb2bbbf8e769ec3da&

From Volfi7
I've gathered all the information I can (if I missed anything, please correct me and I'll add it here) to avoid spamming with recurring bugs.

Current known bugs:
Server disconnection when another player leaves (solution: logging out not close app after entering rebuild mode and waiting a few seconds)
Items in the fridge/gift box being duplicated after a crash 
Food being sold immediately after a crash not all, but a large portion
The wheel in the market doesn't provide decorations (the wheel that pops up in the cafe gives decorations, but you still have to pay for them)
Walls and other decorations aren't saved after a disconnection (solution: by placing walls/decorations every few decorations, exit the server [see point 1])
Counters placed one square apart from the chair cause the customer to receive ghost food, which will give you negative popularity
A counter placed near the door that waiters will use immediately after entering, carrying ghost food, has the same effect as above (solution: move the counter                 so the gap is larger). (larger than 1 square)
Quick cooking sometimes works, sometimes doesn't
You can join a CO-OP, but you can't leave it
After entering the market, you can't visit another player using the ranking
Giving out food from a gift box has a few problems: it's not always sold, and after a crash, it returns to the gift box, leaving the counter glitched (either it's "free"                but visually there's still food on it, or the number of a given food visually increases to the recipe, but the waiter won't give it out [solution: remove the                         food from the counter])
Removing food/spices from the gift box doesn't work
The wheel gives empty winnings
The courier doesn't work
Rearranging items in the cafe has a strange effect on luxury points

Unimplemented things:
Expanding the cafe
Smoothie machine (u can buy, but not use)
Muffin Man
Working for others
Adding to friends

I cannot clear my Gifts tab it is stuck on 94 or 93 every single time and when I relog into the game and I had cleared most of it out it goes back to 94 every time
https://cdn.discordapp.com/attachments/1453807578936115414/1455525932461719605/image.png?ex=6957aea8&is=69565d28&hm=3677b994fe9a2a5fb90605b171a8c1c34e8edc078ec2d036196ce68de898ba99&

I did the Spaghetti Bolognese Coop and payed 1 Gold for it to go 10h longer. After that I would have gotten the gold Time, but only got bronze

charakter getting stuck trying to bring food to a counter that is not reachable
https://cdn.discordapp.com/attachments/1453807578936115414/1455969263327383804/image.png?ex=6957514a&is=6955ffca&hm=c4182332d43987707c1d57044f3dac43efb3cc5dd4fd10c95198ac64fd65d851&

I accidentally cooked the wrong dish, Chocolate fondue that takes 2.5 hours, and threw it away immediately to swap it for the Black Forest cake (13 hours), and now the cooking time for the cake is only 2.5 hours. Several crashes happened in between.
https://cdn.discordapp.com/attachments/1453807578936115414/1456087571234033675/image.png?ex=6957bf79&is=69566df9&hm=cc5c3963993f593196e727597039c52b1e3ba0c66b0bdbb297591a2298c7b9af&
https://cdn.discordapp.com/attachments/1453807578936115414/1456087577462440169/image.png?ex=6957bf7b&is=69566dfb&hm=770ed5988ddd0343f8c36fc9ea35939821ca57d2d98d8921e0319694c03b3480&
