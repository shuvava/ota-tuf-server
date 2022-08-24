package data

import "time"

func validExpires(expires time.Time) bool {
	return time.Until(expires) > 0
}

/*
type SignedPayload struct {
	// TODO (vs): not sure about field type
	Signed     interface{}                  `json:"signed"`
	Signatures []encryption.ClientSignature `json:"signatures"`
}
*/
/*
retry_command "reposerver" "[[ true = \$(http --print=b GET ${reposerver}/health/dependencies \
| jq --exit-status '.status == \"OK\"') ]]"

http_2xx_or_4xx --ignore-stdin --check-status POST "${reposerver}/api/v1/user_repo" "${namespace}"

sleep 5s
id=$(http --ignore-stdin --check-status --print=h GET "${reposerver}/api/v1/user_repo/root.json" "${namespace}" \
| grep -i x-ats-tuf-repo-id | awk '{print $2}' | tr -d '\r')

http_2xx_or_4xx --ignore-stdin --check-status POST "${director}/api/v1/admin/repo" "${namespace}"

retry_command "keys" "http --ignore-stdin --check-status GET ${keyserver}/api/v1/root/${id}"
keys=$(http --ignore-stdin --check-status GET "${keyserver}/api/v1/root/${id}/keys/targets/pairs")
echo ${keys} | jq '.[0] | {keytype, keyval: {public: .keyval.public}}'   > "${SERVER_DIR}/targets.pub"
echo ${keys} | jq '.[0] | {keytype, keyval: {private: .keyval.private}}' > "${SERVER_DIR}/targets.sec"

retry_command "root.json" "http --ignore-stdin --check-status -d GET \
${reposerver}/api/v1/user_repo/root.json \"${namespace}\"" && \
	http --ignore-stdin --check-status -d -o "${SERVER_DIR}/root.json" GET \
		 ${reposerver}/api/v1/user_repo/root.json "${namespace}"

*/

//// SerializedKey is TUF struct for signature and encryption keys
//// in ota-tuf project it is named `TufKey`
//type SerializedKey struct {
//	// Type is key type
//	Type KeyType `json:"keytype"`
//	// Value is key value
//	Value json.RawMessage `json:"keyval"`
//}
//
//// PrivateKey is a private key
//type PrivateKey struct {
//	SerializedKey
//}
//
//// PublicKey is a public key
//type PublicKey struct {
//	SerializedKey
//	ids []string
//}

//// TufKeyType is named as KeyType in original code
//// RsaKeyType,Ed25519KeyType,EcPrime256KeyType types skipped for now
//type TufKeyType int
//
//const (
//	Pub TufKeyType = iota
//	Priv
//	Pair
//)

//type TufKey struct {
//	ID      KeyID      `json:"id"`
//	KeyType TufKeyType `json:"keytype"`
//	KeyVal  PublicKey  `json:"keyval"`
//}
//
//func RSATufKey(keyval PublicKey) TufKey {
//	return TufKey{
//		ID:      NewKeyID(keyval),
//		KeyType: Pub,
//		KeyVal:  keyval,
//	}
//}
