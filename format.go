package main

import "fmt"

func fmtPrice(p int) string {
	cents := p % 100
	euros := int((p - cents) / 100)
	return fmt.Sprintf("%d,%d", euros, cents)
}

func getMessage(f []DBEntry, t string) string {

	foods := ""

	for _, entry := range f {
		if len(entry.Items) == 0 {
			continue
		}

		foods += fmt.Sprintf(`,{"type": "section","text": {"type": "mrkdwn","text": "*%s*`, entry.Canteen)

		for _, item := range entry.Items {
			foods += "\n"

			if !entry.SpecDiet {
				foods += ":black_small_square:"
			} else if item.Vegan {
				foods += ":seedling:"
			} else if item.Vegetarian {
				foods += ":cheese_wedge:"
			} else {
				foods += ":cut_of_meat:"
			}

			foods += item.Name
			foods += " € "
			foods += fmtPrice(item.StudPrice)

			if item.StudPrice != item.ProfPrice {
				foods += "/" + fmtPrice(item.ProfPrice)
			}

		}
		foods += `"}}`
	}

	return fmt.Sprintf(`{"blocks": [{"type": "section","text": {"type": "mrkdwn","text": "%s"}},{"type": "divider"}%s]}`, t, foods)
}