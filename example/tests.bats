#!/usr/bin/env bats

load '/bats-lib/bats-support/load'
load '/bats-lib/bats-assert/load'

@test "test without mtls succeeds" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           https://whoami.7f000001.nip.io/
  assert_success
  assert_output --partial 'X-Forwarded-Server'
}

@test "test mtls without cert is rejected" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           https://whoami-mtls.7f000001.nip.io/
  assert_failure 56
  assert_output --partial 'alert bad certificate, errno 0'
}

@test "test mtls with bad-client cert is rejected" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/bad-client/cert.pem \
           --key traefik_config/certs/bad-client/key.pem \
           https://whoami-mtls.7f000001.nip.io/
  assert_failure 56
  assert_output --partial 'alert bad certificate, errno 0'
}

@test "test mtls with auth-client cert is accepted" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/auth-client/cert.pem \
           --key traefik_config/certs/auth-client/key.pem \
           https://whoami-mtls.7f000001.nip.io/
  assert_success
  assert_output --partial 'X-Forwarded-Server'
}

@test "test mtls with auth2-client cert is also accepted" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/auth2-client/cert.pem \
           --key traefik_config/certs/auth2-client/key.pem \
           https://whoami-mtls.7f000001.nip.io/
  assert_success
  assert_output --partial 'X-Forwarded-Server'
}
