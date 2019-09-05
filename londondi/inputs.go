package londondi

import (
	"fmt"
	"strconv"

	"github.com/byuoitav/common/pooled"
)

//PORT .
const PORT = "1023"

//SetVolume gets an volume level (int) and sets it on device
func SetVolume(input string, address string, volume string) error {
	_, err := strconv.Atoi(volume)
	if err != nil {
		return err
	}

	command, err := BuildRawVolumeCommand(input, address, volume)
	if err != nil {
		return err
	}

	command, err = MakeSubstitutions(command, ENCODE)
	if err != nil {
		return err
	}

	command, err = Wrap(command)
	if err != nil {
		return err
	}

	work := func(conn pooled.Conn) error {
		n, err := conn.Write(command)
		switch {
		case err != nil:
			return fmt.Errorf("unable to send command: %s", err)
		case n != len(command):
			return fmt.Errorf("unable to send command: wrote %v/%v bytes", n, len(command))
		}
		return nil
	}

	err = pool.Do(address, work)
	if err != nil {
		return err
	}

	return nil
}

//UnMute handler
func UnMute(input string, address string) error {
	command, err := BuildRawMuteCommand(input, address, "false")
	if err != nil {
		return err
	}

	command, err = MakeSubstitutions(command, ENCODE)
	if err != nil {
		return err
	}

	command, err = Wrap(command)
	if err != nil {
		return err
	}

	work := func(conn pooled.Conn) error {
		n, err := conn.Write(command)
		switch {
		case err != nil:
			return fmt.Errorf("unable to send command: %s", err)
		case n != len(command):
			return fmt.Errorf("unable to send command: wrote %v/%v bytes", n, len(command))
		}
		return nil
	}

	err = pool.Do(address, work)
	if err != nil {
		return err
	}

	return nil
}

//Mute .
func Mute(input string, address string) error {
	command, err := BuildRawMuteCommand(input, address, "true")
	if err != nil {
		return err
	}

	command, err = MakeSubstitutions(command, ENCODE)
	if err != nil {
		return err
	}

	command, err = Wrap(command)
	if err != nil {
		return err
	}

	work := func(conn pooled.Conn) error {
		n, err := conn.Write(command)
		switch {
		case err != nil:
			return fmt.Errorf("unable to send command: %s", err)
		case n != len(command):
			return fmt.Errorf("unable to send command: wrote %v/%v bytes", n, len(command))
		}
		return nil
	}

	err = pool.Do(address, work)
	if err != nil {
		return err
	}

	return nil
}
