package cleanrooms

import (
	"context"
	"errors"

	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"
	"github.com/aws/aws-sdk-go-v2/service/cleanrooms/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"

	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_cleanrooms_configured_table_association")
// @Tags(identifierAttribute="arn")
func ResourceConfiguredTableAssociation() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceConfiguredTableAssociationCreate,
		ReadWithoutTimeout:   resourceConfiguredTableAssociationRead,
		UpdateWithoutTimeout: resourceConfiguredTableAssociationUpdate,
		DeleteWithoutTimeout: resourceConfiguredTableAssociationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrName: {
				Type:     schema.TypeString,
				Required: true,
			},
			"membership_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configured_table_identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameConfiguredTableAssociation = "Configured Table Association"
)

func resourceConfiguredTableAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CleanRoomsClient(ctx)

	input := &cleanrooms.CreateConfiguredTableAssociationInput{
		Name:                      aws.String(d.Get(names.AttrName).(string)),
		MembershipIdentifier:      aws.String(d.Get("membership_identifier").(string)),
		ConfiguredTableIdentifier: aws.String(d.Get("configured_table_identifier").(string)),
		RoleArn:                   aws.String(d.Get("role_arn").(string)),
		Tags:                      getTagsIn(ctx),
	}

	if v, ok := d.GetOk(names.AttrDescription); ok {
		input.Description = aws.String(v.(string))
	}

	out, err := conn.CreateConfiguredTableAssociation(ctx, input)

	if err != nil {
		return create.DiagError(names.CleanRooms, create.ErrActionCreating, ResNameConfiguredTableAssociation, d.Get("name").(string), err)
	}

	if out == nil || out.ConfiguredTableAssociation == nil {
		return create.DiagError(names.CleanRooms, create.ErrActionCreating, ResNameConfiguredTableAssociation, d.Get("name").(string), errors.New("empty output"))
	}

	d.SetId(aws.ToString(out.ConfiguredTableAssociation.Id))
	return nil
}

func resourceConfiguredTableAssociationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CleanRoomsClient(ctx)

	out, err := findConfiguredTableAssociationByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CleanRooms Configured Table Association (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.CleanRooms, create.ErrActionReading, ResNameConfiguredTableAssociation, d.Id(), err)
	}

	configuredTableAssociation := out.ConfiguredTableAssociation
	d.Set(names.AttrARN, configuredTableAssociation.Arn)
	d.Set(names.AttrName, configuredTableAssociation.Name)
	d.Set(names.AttrDescription, configuredTableAssociation.Description)
	d.Set("create_time", configuredTableAssociation.CreateTime.String())
	d.Set("update_time", configuredTableAssociation.UpdateTime.String())
	d.Set("membership_identifier", configuredTableAssociation.MembershipId)

	return nil
}

func resourceConfiguredTableAssociationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CleanRoomsClient(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &cleanrooms.UpdateConfiguredTableAssociationInput{
			ConfiguredTableAssociationIdentifier: aws.String(d.Id()),
		}

		if d.HasChanges(names.AttrDescription) {
			input.Description = aws.String(d.Get(names.AttrDescription).(string))
		}

		_, err := conn.UpdateConfiguredTableAssociation(ctx, input)
		if err != nil {
			return create.DiagError(names.CleanRooms, create.ErrActionUpdating, ResNameConfiguredTableAssociation, d.Id(), err)
		}
	}

	return append(diags, resourceConfiguredTableAssociationRead(ctx, d, meta)...)
}

func resourceConfiguredTableAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CleanRoomsClient(ctx)

	log.Printf("[INFO] Deleting CleanRooms Table Association %s", d.Id())
	_, err := conn.DeleteConfiguredTableAssociation(ctx, &cleanrooms.DeleteConfiguredTableAssociationInput{
		ConfiguredTableAssociationIdentifier: aws.String(d.Id()),
	})

	if errs.IsA[*types.AccessDeniedException](err) {
		return nil
	}

	if err != nil {
		return create.DiagError(names.CleanRooms, create.ErrActionDeleting, ResNameConfiguredTableAssociation, d.Id(), err)
	}

	return nil
}
