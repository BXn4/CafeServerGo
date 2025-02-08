
# NPC Design
This file contains the details of the npc design

## Waiters
The waiters spawn every time the player enters and they dont despawn until it comes back
Each waiter needs a different switch so we can fire them
Because of these we need full control of waiters

The waiter.go implements a agent with a finite state machine with 4 states:
- getAndMoveToCounter:
  This is the first state
  In this state we choose a random counter we prioritize counters with dishes then move to it if we cant get a counter we stop
- selectJob:
  - In this state we roll a random number and choose a job based on waiter priority
- takePlates
  - In this state
- serveFood:

![Waiter states](../../assets/imgs/waiter_states.svg "Waiter states")

## Customers
