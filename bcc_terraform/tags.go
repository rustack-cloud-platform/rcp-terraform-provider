package bcc_terraform

import (
	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func marshalTagNames(tags []bcc.Tag) []interface{} {
	convertedTags := make([]interface{}, len(tags))
	for i, tag := range tags {
		convertedTags[i] = tag.Name
	}
	return convertedTags
}

func unmarshalTagNames(tags interface{}) []bcc.Tag {
	tagList := tags.(*schema.Set).List()
	resultTags := make([]bcc.Tag, len(tagList))
	for i, tag := range tagList {
		resultTags[i] = bcc.Tag{Name: tag.(string)}
	}
	return resultTags
}

func newTagNamesResourceSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:        schema.TypeString,
			Description: "name of the Tag",
		},
		Description: description,
	}
}

func marshalTags(tags []bcc.Tag) []map[string]interface{} {
	convertedTags := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		convertedTags[i]["id"] = map[string]interface{}{"id": tag.ID, "name": tag.Name}
	}
	return convertedTags
}

func newTagsDatasourceSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: Arguments{
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "id of the Tag",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "name of the Tag",
				},
			},
		},
		Description: description,
	}
}
