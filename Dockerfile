FROM golang:1.16.3-alpine
# アップデートとgitのインストール
RUN apk update && apk add git && apk add build-base
# appディレクトリの作成
RUN mkdir /omusubi-clone
RUN go get github.com/labstack/echo/v4
RUN go get github.com/labstack/echo/v4/middleware
RUN go get -u -v bitbucket.org/liamstask/goose/cmd/goose
# ホットリロード追加
RUN go get -u -v github.com/pilu/fresh
# RUN go get github.com/pressly/goose
#ワーキングディレクトリの設定
WORKDIR /omusubi-clone
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /omusubi-clone