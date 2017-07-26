/**
 * caddy-docker Web UI
 * Copyright 2017 - Nathan Osman
 */

(function() {

    /**
     * Create a <td> element
     */
    function td(t) {
        if ($.isArray(t)) {
            t = t.join(', ');
        }
        return $('<td>').text(t);
    }

    /**
     * Reloads the list of currently running containers
     */
    function reloadContainers() {
        $.post('/api', {action: 'getContainers'}, function(d) {

            // Sort the containers by name
            d = d.sort(function(a, b) {
                a = a.Name.toUpperCase();
                b = b.Name.toUpperCase();
                if (a < b) {
                    return -1;
                } if (a > b) {
                    return 1;
                } else {
                    return 0;
                }
            });

            // Add the containers to the list
            var $tbody = $('#containers tbody').empty();
            $.each(d, function() {
                $('<tr>')
                    .append(td(this.Name))
                    .append(td(this.Domains))
                    .append(td(this.Addr))
                    .appendTo($tbody);
            });
        });
    }

    reloadContainers();

})();
