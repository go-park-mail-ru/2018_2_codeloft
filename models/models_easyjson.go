// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user_id":
			out.Id = int64(in.Int64())
		case "login":
			out.Login = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "score":
			out.Score = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Id))
	}
	{
		const prefix string = ",\"login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"score\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Score))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels1(in *jlexer.Lexer, out *Leaders) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make([]User, 0, 1)
					} else {
						out.Users = []User{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v1 User
					if data := in.Raw(); in.Ok() {
						in.AddError((v1).UnmarshalJSON(data))
					}
					out.Users = append(out.Users, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "pagesCount":
			out.PagesCount = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels1(out *jwriter.Writer, in Leaders) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"users\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Users {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.Raw((v3).MarshalJSON())
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"pagesCount\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.PagesCount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Leaders) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Leaders) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Leaders) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Leaders) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels1(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels2(in *jlexer.Lexer, out *HelpUser) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "new_password":
			out.NewPassword = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "score":
			out.Score = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels2(out *jwriter.Writer, in HelpUser) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	if in.NewPassword != "" {
		const prefix string = ",\"new_password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.NewPassword))
	}
	if in.Email != "" {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	if in.Score != 0 {
		const prefix string = ",\"score\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Score))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v HelpUser) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v HelpUser) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20182CodeloftModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *HelpUser) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *HelpUser) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20182CodeloftModels2(l, v)
}
