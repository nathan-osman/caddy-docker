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
