# Dcard Homework - middleware

## 題目

Dcard 每天午夜都有大量使用者湧入抽卡，為了不讓伺服器過載，請設計一個 middleware：

- 限制每小時來自同一個 IP 的請求數量不得超過 1000
- 在 response headers 中加入剩餘的請求數量 (X-RateLimit-Remaining) 以及 rate limit 歸零的時間 (X-RateLimit-Reset)
- 如果超過限制的話就回傳 429 (Too Many Requests)
- 可以使用各種資料庫達成

## 想法

因申請條件有一條是熟悉Golang或node.js的框架，剛好我現在主力語言是go，雖然沒用go寫過web，但有機會練習也不錯

資料庫本來想用PostgreSQL，後來想想覺得用redis更好，因為
它是單純的key-value，還有天然內建的TTL可以使用

趁這個機會同時練習go的http和redis操作，挺不錯的

## TODO

- ~~fix bugs~~
- ~~test~~
- ~~CI~~
- log
- framework
