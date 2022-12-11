# Ding

> A low-effort CLI-based social network

## User Documentation

Listen for dings as alice@ding.example.net: `ding listen alice@ding.example.net`

Send a ding to alice@ding.example.net: `ding ding alice@ding.example.net`

ding.example.net is a ding server (see below)

dings arrive as a terminal bell, so you need to have your terminal configured to ring the bell correctly

## Server Documentation

You can run your own ding server quite easily: `ding server`

The only requirement is that TCP port 1883 is open

## FAQs

* *Can I send anything other than a ding?* No
* *What about authentication?* There is none
* *What if someone else pretends to be me?* You will be disconnected
* *What if a ding is sent while I am not listening?* It will be lost forever
* *What does a ding mean?* Whatever you want it to mean