# markdown-parser-go  
Markdown形式のテキストをHTML形式へ変換します。  
GitHub Flavored MarkdownやGitLab Flavored Markdownなど様々な方言が存在しますが、独自の方言を実装します。  
そのため、他の方言と書き方が異なる場合があります。 

## 使用技術
- Go 1.20.2

## 対応状況
- [x] 強調(strong)
- [x] 箇条書きリスト(ul)
- [x] 順序付きリスト(ol)
- [x] 見出し
- [x] リンク

## 使用方法
### 1. dockerコンテナの起動  
以下のコマンドを実行するとdockerコンテナが起動します。  
```bash
docker-compose up
```

### 2. dockerコンテナに接続  
以下のコマンドを実行するとmarkdown-parser-goコンテナに接続します。  
```bash
docker exec -it markdown-parser-go bash
```

### 3. サーバーを起動する  
以下のコマンドを実行するとサーバーが起動します。  
```bash
go run main.go
```

### 4. リクエストを送信する  
以下のコマンドを実行することで{markdown_text}をサーバーに送信することができます。  
リクエスト送信後、HTML形式に変換された{markdown_text}が表示されます。  
コマンドの実行が上手くいかない場合には、ターミナルを別で開き、コンテナに接続してコマンドを実行してください。  
POSTmanなどのアプリを使ってリクエストを送信することも可能です。  
```bash
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"Content": "{markdown text}"}' http://localhost:8081/convertmd
```

## 参考サイト
- [マークダウンパーサを作ろう - エムスリーテックブログ](https://www.m3tech.blog/entry/2021/08/23/124000)