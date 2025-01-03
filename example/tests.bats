#!/usr/bin/env bats

load '/bats-lib/bats-support/load'
load '/bats-lib/bats-assert/load'

@test "test whoami is working" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           https://whoami.7f000001.nip.io/
  assert_success
  assert_output --partial 'X-Forwarded-Server'
}

@test "test $backend without cert is rejected" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           https://$backend.7f000001.nip.io/
  assert_failure 56
  assert_output --partial 'alert bad certificate, errno 0'
}

@test "test $backend with bad-client cert is rejected" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/bad-client/cert.pem \
           --key traefik_config/certs/bad-client/key.pem \
           https://$backend.7f000001.nip.io/
  assert_failure 56
  assert_output --partial 'alert bad certificate, errno 0'
}

@test "test $backend with auth-client cert is accepted" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/auth-client/cert.pem \
           --key traefik_config/certs/auth-client/key.pem \
           https://$backend.7f000001.nip.io/
  assert_success
  assert_output --partial 'X-Forwarded-Server'
}

@test "test $backend with auth2-client cert is also accepted" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/auth2-client/cert.pem \
           --key traefik_config/certs/auth2-client/key.pem \
           https://$backend.7f000001.nip.io/
  assert_success
  assert_output --partial 'X-Forwarded-Server'
}

@test "test $backend with not-auth-client cert is denied" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/not-auth-client/cert.pem \
           --key traefik_config/certs/not-auth-client/key.pem \
           https://$backend.7f000001.nip.io/
  assert_failure 22
  assert_output --partial 'returned error: 403'
}
