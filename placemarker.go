package placemarker

import (
	"fmt"
	"os"

	"github.com/twpayne/go-kml"
)

func AddPoint(k *kml.CompoundElement, name, desc string, alt, lat, lon float64) {
	k.Add(
		kml.Placemark(
			kml.Name(name),
			kml.Description(desc),
			kml.Point(
				kml.Coordinates(
					kml.Coordinate{
						Lon: lon,
						Lat: lat,
						Alt: alt,
					}),
			),
		))
}

func WriteKML(k *kml.CompoundElement, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	if err := k.Write(f); err != nil {
		return fmt.Errorf("k.Write: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	return nil
}
