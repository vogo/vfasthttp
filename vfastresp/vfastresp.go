/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vfastresp

import (
	"encoding/json"
	"html/template"
	"log"

	"github.com/vogo/vogo/vnet/vhttp/vhttperror"
	"github.com/vogo/vogo/vnet/vhttp/vhttpresp"

	"github.com/valyala/fasthttp"
)

func Data(ctx *fasthttp.RequestCtx, code int, data interface{}) {
	Write(ctx, code, "", data)
}

func CodeData(ctx *fasthttp.RequestCtx, code int, msg string, data interface{}) {
	Write(ctx, code, msg, data)
}

func OK(ctx *fasthttp.RequestCtx) {
	Write(ctx, vhttperror.CodeOK, "", "ok")
}

func Success(ctx *fasthttp.RequestCtx, data interface{}) {
	Write(ctx, vhttperror.CodeOK, "", data)
}

func CodeError(ctx *fasthttp.RequestCtx, code int, err error) {
	CodeMsg(ctx, code, err.Error())
}

func Error(ctx *fasthttp.RequestCtx, err error) {
	if c, ok := err.(vhttperror.StatusState); ok {
		ctx.SetStatusCode(c.Status())
	}

	code := vhttperror.CodeUnknownErr

	if c, ok := err.(vhttperror.Coder); ok {
		code = c.Code()
	}

	CodeMsg(ctx, code, err.Error())
}

func BadMsg(ctx *fasthttp.RequestCtx, msg string) {
	CodeMsg(ctx, vhttperror.CodeBadRequestErr, msg)
}

func BadError(ctx *fasthttp.RequestCtx, err error) {
	BadMsg(ctx, err.Error())
}

func CodeMsg(ctx *fasthttp.RequestCtx, code int, msg string) {
	Write(ctx, code, msg, nil)
}

func Write(ctx *fasthttp.RequestCtx, code int, msg string, data interface{}) {
	resp := vhttpresp.ResponseBody{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json marshal error: %+v", err)

		_, _ = ctx.Write([]byte("internal server error"))

		return
	}

	ctx.SetContentType("application/json")
	_, _ = ctx.Write(jsonBytes)
}

func Template(ctx *fasthttp.RequestCtx, tpl *template.Template, data interface{}) {
	ctx.SetContentType("text/html")
	err := tpl.Execute(ctx.Response.BodyWriter(), data)
	if err != nil {
		log.Fatalf("template format error: %v", err)
	}
}
