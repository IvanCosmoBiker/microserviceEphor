package command

import (
)

type CommandRequest struct {
    D string 
    A int
}

func (c *CommandRequest) GetDeviceNumber() string {
    return c.D
}

func (c *CommandRequest) GetActionOfDevice() int {
    return c.A
}

type CommandResponse struct {
    A int64
    Sum int64
    M int
}

type CommandServerRequest struct {
    Id string
    Imei string
}


