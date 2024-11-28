#!/usr/bin/env bats

@test "wrapper test" {
    final_status=0
    BACKENDS=("whoami-mtls" "whoami-mtls-raw")

    for b in "${BACKENDS[@]}"; do
        backend="$b" run bats -t tests.bats
        echo "# $output" >&3
        echo "#" >&3
        final_status=$(($final_status + $status))
    done

    [ "$final_status" -eq 0 ]
}
