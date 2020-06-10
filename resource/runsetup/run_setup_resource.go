package runsetup

import (
	"github.com/draganm/terraform-provider-linuxbox/sshsession"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreate,
		Read:   resourceRead,
		Update: resourceUpdate,
		Delete: resourceDelete,

		Schema: map[string]*schema.Schema{
			"ssh_key": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"ssh_user": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "root",
				Optional: true,
			},

			"host_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"setup": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},

			"check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"delete": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCreate(d *schema.ResourceData, m interface{}) error {

	ssh, err := sshsession.Open(d)
	if err != nil {
		return errors.Wrap(err, "while creating ssh session")
	}

	defer ssh.Close()

	setup := d.Get("setup").([]interface{})

	for _, sl := range setup {
		line := sl.(string)
		stdout, stderr, err := ssh.RunInSession(line)
		if err != nil {
			return errors.Wrapf(err, "error while executing %q\nSTDOUT:\n%s\nSTDERR:\n%s\n", line, string(stdout), string(stderr))
		}
	}

	d.SetId("-")

	return nil
}

func resourceRead(d *schema.ResourceData, m interface{}) error {

	ssh, err := sshsession.Open(d)
	if err != nil {
		return errors.Wrap(err, "while creating ssh session")
	}

	defer ssh.Close()

	check, checkSet := d.GetOkExists("check")

	if !checkSet {
		return nil
	}

	_, _, err = ssh.RunInSession(check.(string))
	if sshsession.IsExecError(err) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return errors.Wrapf(err, "while running check")
	}

	return nil
}

func resourceUpdate(d *schema.ResourceData, m interface{}) error {
	return errors.New("update is not supported")
}

func resourceDelete(d *schema.ResourceData, m interface{}) error {
	ssh, err := sshsession.Open(d)
	if err != nil {
		return errors.Wrap(err, "while creating ssh session")
	}

	defer ssh.Close()

	delete, deleteSet := d.GetOkExists("delete")

	if !deleteSet {
		return nil
	}

	stdout, stderr, err := ssh.RunInSession(delete.(string))
	if err != nil {
		return errors.Wrapf(err, "error while executing %q\nSTDOUT:\n%s\nSTDERR:\n%s\n", delete, string(stdout), string(stderr))
	}

	return nil
}
