// Package saml knows how to interact with a SAML assertion document
package saml

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

const errMsg = `an error occurred when processing the SAML response from AWS:
%w
`

// nolint: stylecheck
func formatErr(err error) error {
	return fmt.Errorf(errMsg, err)
}

// VerifyAssertion tries to find specific information we are looking for
// Much of this is borrowed from:
// - https://github.com/Versent/saml2aws/blob/master/saml.go
func VerifyAssertion(expectRole string, assertion []byte) error {
	data, err := base64.StdEncoding.DecodeString(string(assertion))
	if err != nil {
		return formatErr(fmt.Errorf("base64 decoding SAML assertion: %w", err))
	}

	doc := etree.NewDocument()

	err = doc.ReadFromBytes(data)
	if err != nil {
		return formatErr(fmt.Errorf("parsing SAML assertion document: %w", err))
	}

	assertionElement := doc.FindElement(".//Assertion")
	if assertionElement == nil {
		return formatErr(fmt.Errorf("missing Assertion tag in saml document"))
	}

	attributeStatement := assertionElement.FindElement(childPath(assertionElement.Space, "AttributeStatement"))
	if attributeStatement == nil {
		return formatErr(fmt.Errorf("missing element AttributeStatement"))
	}

	foundRole := false

	attributes := attributeStatement.FindElements(childPath(assertionElement.Space, "Attribute"))
FOUNDROLE:
	for _, attribute := range attributes {
		if attribute.SelectAttrValue("Name", "") != "https://aws.amazon.com/SAML/Attributes/Role" {
			continue
		}
		attributeValues := attribute.FindElements(childPath(assertionElement.Space, "AttributeValue"))
		for _, attrValue := range attributeValues {
			if strings.Contains(attrValue.Text(), expectRole) {
				foundRole = true
				break FOUNDROLE
			}
		}
	}

	if !foundRole {
		return formatErr(fmt.Errorf("you do not have permission to use the role: %s, ask for help in #kjørermiljø-support on slack", expectRole))
	}

	return nil
}

func childPath(space, tag string) string {
	if space == "" {
		return "./" + tag
	}

	return "./" + space + ":" + tag
}
