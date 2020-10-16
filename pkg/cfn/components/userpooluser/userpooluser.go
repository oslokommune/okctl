// Package userpooldomain provides functionality for setting
// up a domain with a user pool
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html
package userpooluser

import (
	"fmt"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// UserPoolUser
type UserPoolUser struct {
	StoredName  string
	Email       string
	Description string
	UserPoolId  string
}


// Well this is new and does not change anything
func (user *UserPoolUser) NamedOutputs() map[string]map[string]interface{} {

	return cfn.NewValue(fmt.Sprintf("%sID", user.StoredName), user.Ref()).NamedOutputs()

	//return cfn.NewValue(fmt.Sprintf("%sClientID", c.Purpose), c.Ref()).NamedOutputs()

}


/*


   mfaRequired: (challengeName, challengeParam) => {
       logger.debug('signIn MFA required');
       user['challengeName'] = challengeName;
       user['challengeParam'] = challengeParam;
       resolve(user);
   },
   newPasswordRequired: (userAttributes, requiredAttributes) => {
       logger.debug('signIn new password');
       user['challengeName'] = 'NEW_PASSWORD_REQUIRED';
       user['challengeParam'] = {
           userAttributes: userAttributes,
           requiredAttributes: requiredAttributes
       };
       resolve(user);
   }
 */



// TODO fix hardcoded group name.
// TODO make linter happy
func (user *UserPoolUser) Resource() cloudformation.Resource {
	u := &cognito.UserPoolUser{
		ClientMetadata:         nil,
		DesiredDeliveryMediums: []string{"EMAIL"},
		ForceAliasCreation:     true,
		MessageAction:          "SUPPRESS",
		UserAttributes: []cognito.UserPoolUser_AttributeType{
			{
				Name:  "email",
				Value: user.Email,
			},
			{
				Name:  "name",
				Value: user.Email,
			},
			{
				Name:  "email_verified",
				Value: "True",
			},
		},
		Username: user.Email,
		UserPoolId:                           user.UserPoolId,
		AWSCloudFormationDeletionPolicy:      "",
		AWSCloudFormationUpdateReplacePolicy: "",
		AWSCloudFormationDependsOn:           nil,
		AWSCloudFormationMetadata:            nil,
		AWSCloudFormationCondition:           "",
	}

	fmt.Println("USER ÆÆÆÆÆÆ")
	fmt.Println(u)


	return u
}

// TODO make linter happy (Something here?)
func (d *UserPoolUser) Name() string {
	return d.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (u *UserPoolUser) Ref() string {
	return cloudformation.Ref(u.Name())
}


// NOW we just need to find where the *fk we actually compose the thing
// TODO make linter happy -- here we go, slug make is non alphanumeric?
//StoredName:  fmt.Sprintf("UserPoolUser%s", strings.ReplaceAll(slug.Make(email), "-", "")),
//StoredName:  fmt.Sprintf("UserPoolUser-%s", slug.Make(email)),

// Could it be something with the username not beeing the email ... also?

func New(email string, description string, userPoolid string) *UserPoolUser {
	// EH --- can it be that we need to match up ... the id?
	return &UserPoolUser{
		StoredName:  "UserPoolUser",
		Email:       email,
		Description: description,
		UserPoolId:  userPoolid,
	}
}
