package main

import (
	"log"
	"os"
)

type AudioPlayer interface {
	Load(*os.File) error
	Play() error
	Pause() error
}

type StringAudioPlayer struct{}

func (p *StringAudioPlayer) Load(f *os.File) error {
	log.Println("Loading")
	return nil
}
func (p *StringAudioPlayer) Play() error {
	log.Println("Play")
	return nil
}
func (p *StringAudioPlayer) Pause() error {
	log.Println("Pause")
	return nil
}

type EventLoad struct {
	File *os.File
}

type EventPlay struct {
}

type EventPause struct {
}

type EventError struct {
	Message string
}

type AudioPlayerState struct {
	file    *os.File
	Player  AudioPlayer
	Message string
}
